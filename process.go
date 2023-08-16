package main

import (
	"bufio"
	"encoding/base64"
	"strings"

	"subruleset/tlog"
)

// processRules 找到以RULE-SET开头，并且第2段以_ipcidr结尾的行，并在这一行最后面加入,no-solve
// 对 Surge 的规则，就修改url的值，变成可读的
func processRules(body, hostname, fullURL string) string {
	tlog.Info.Println("Processing rules")

	scanner := bufio.NewScanner(strings.NewReader(body))
	var result []string

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#!MANAGED-CONFIG") {
			parts := strings.Split(line, " ")
			if len(parts) > 1 {
				line = strings.Replace(line, parts[1], fullURL, 1)
			}
		} else if strings.HasPrefix(line, "RULE-SET,https://sub.zaptiah.com/getruleset?type=1") {
			line = processSurgeRule(line, hostname)
		} else if strings.HasPrefix(line, "  - RULE-SET,") && strings.HasSuffix(strings.Split(line, ",")[1], "_ipcidr") {
			line = line + ",no-resolve"
		}
		result = append(result, line)
	}

	tlog.Info.Println("Finished processing rules")
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
