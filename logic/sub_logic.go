package logic

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"subruleset/tlog"
)

type SubLogic struct {
	url   string
	token string
}

func NewSubLogic(url, token string) *SubLogic {
	return &SubLogic{
		url:   url,
		token: token,
	}
}

// ValidateToken 验证token
func (l *SubLogic) ValidateToken(token string) bool {
	return token == l.token
}

// FetchSubscriptions 从机场获取订阅
func (l *SubLogic) FetchSubscriptions() (http.Header, string, error) {
	resp, err := http.Get(l.url)
	if err != nil {
		return nil, "", fmt.Errorf("error fetching data from URL: %s, err: %v", l.url, err)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("error reading response body: %v", err)
	}
	return resp.Header, string(bytes), nil
}

// Clash 处理Clash规则，添加no-resolve
func (l *SubLogic) Clash(body string) string {
	tlog.Info.Println("Processing Clash rules")

	scanner := bufio.NewScanner(strings.NewReader(body))
	var result []string

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "  - RULE-SET,") && strings.HasSuffix(strings.Split(line, ",")[1], "_ipcidr") {
			line = line + ",no-resolve"
		}
		result = append(result, line)
	}

	tlog.Info.Println("Finished processing Clash rules")
	return strings.Join(result, "\n")
}

// Surge 处理Surge规则
// 1. 修改订阅地址
// 2. 使RULE-SET的url变成可读
// 3. 添加dns-failed
func (l *SubLogic) Surge(body string, currentURL string) string {
	tlog.Info.Println("Processing Surge rules")

	parsedURL, err := url.Parse(currentURL)
	if err != nil {
		// 处理错误，例如返回错误信息或默认值
		return fmt.Sprintf("Error parsing URL: %v", err)
	}

	hostname := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)

	scanner := bufio.NewScanner(strings.NewReader(body))
	var result []string

	for scanner.Scan() {
		line := scanner.Text()
		// 修改订阅地址
		if strings.HasPrefix(line, "#!MANAGED-CONFIG") {
			parts := strings.Split(line, " ")
			if len(parts) > 1 {
				line = strings.Replace(line, parts[1], currentURL, 1)
			}
		} else if strings.HasPrefix(line, "RULE-SET,https://sub.zaptiah.com/getruleset?type=1") {
			line = processSurgeRule(line, hostname)
		} else if strings.HasPrefix(line, "FINAL,") {
			line = line + ", dns-failed"
		}
		result = append(result, line)
	}

	tlog.Info.Println("Finished processing Surge rules")
	return strings.Join(result, "\n")
}

// processSurgeRule 处理Surge规则，主要把url后面的值变成可读
func processSurgeRule(line, hostname string) string {
	tlog.Info.Println("Processing Surge rule: ", line)

	// Extract the url value
	urlValueStartIdx := strings.Index(line, "url=") + 4
	urlValueEndIdx := strings.Index(line[urlValueStartIdx:], ",")
	if urlValueEndIdx == -1 {
		urlValueEndIdx = len(line)
	} else {
		urlValueEndIdx += urlValueStartIdx
	}

	urlValue := line[urlValueStartIdx:urlValueEndIdx]

	// Base64 decode the value
	decodedValue, err := base64.RawStdEncoding.DecodeString(urlValue)
	if err != nil {
		tlog.Error.Println("Error decoding base64 value: ", err)
		return line
	}

	// Extract the last part of the decoded value
	parts := strings.Split(string(decodedValue), "/")
	lastPart := parts[len(parts)-1]

	// Construct the new line
	typePart := strings.TrimSuffix(lastPart, ".list")
	newLine := "RULE-SET," + hostname + "/rule?type=" + typePart

	// Append other parts from the original line
	remainingParts := strings.SplitN(line[urlValueEndIdx:], ",", 3)
	if len(remainingParts) >= 3 {
		newLine += "," + remainingParts[1] + "," + remainingParts[2]
	}

	return newLine
}