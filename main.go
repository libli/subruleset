package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// 以下为日志级别分级
var (
	Info  = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Warn  = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
)

type Config struct {
	Token string            `yaml:"token"`
	Urls  map[string]string `yaml:"urls"`
}

func main() {
	config, err := readConfig("config.yaml")
	if err != nil {
		Error.Fatalf("Failed to read config: %v", err)
	}

	if config.Token == "" {
		Error.Fatal("Token is empty")
	}

	mux := http.NewServeMux()
	for key, url := range config.Urls {
		Info.Printf("Adding handler for path: %s, url: %s", key, url)
		mux.HandleFunc(fmt.Sprintf("/%s", key), handler(url, config.Token))
	}

	Info.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", mux)
}

// readConfig 读取配置文件
func readConfig(filename string) (Config, error) {
	Info.Printf("Reading configuration from file: %s", filename)

	var config Config
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	Info.Println("Successfully read configuration")
	return config, nil
}

// handler 处理请求
func handler(url string, confToken string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		Info.Printf("Handling request for URL: %s", r.URL)

		token := r.URL.Query().Get("token")
		if token != confToken {
			Warn.Printf("Forbidden request made with token: %s", token)
			http.Error(w, "无权限", http.StatusForbidden)
			return
		}
		resp, err := http.Get(url)
		if err != nil {
			Error.Printf("Error fetching data from URL: %s, err: %v", url, err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// 设置 "access-control-allow-origin" 为 "*"
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// 如果有 "subscription-userinfo" 这个header，那么复制它
		if values, ok := resp.Header["Subscription-Userinfo"]; ok {
			for _, value := range values {
				w.Header().Add("Subscription-Userinfo", value)
			}
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			Error.Printf("Error reading response body: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		newBody := processRules(string(body), r.Host)

		w.Write([]byte(newBody))
	}
}

// processRules 找到以RULE-SET开头，并且第2段以_ipcidr结尾的行，并在这一行最后面加入,no-solve
// 对 Surge 的规则，就修改url的值，变成可读的
func processRules(body, hostname string) string {
	Info.Println("Processing rules")

	scanner := bufio.NewScanner(strings.NewReader(body))
	var result []string

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "RULE-SET,https://sub.zaptiah.com/getruleset?type=1") {
			line = processSurgeRule(line, hostname)
		} else if strings.HasPrefix(line, "  - RULE-SET,") && strings.HasSuffix(strings.Split(line, ",")[1], "_ipcidr") {
			line = line + ",no-resolve"
		}
		result = append(result, line)
	}

	Info.Println("Finished processing rules")
	return strings.Join(result, "\n")
}

// processSurgeRule 处理Surge规则，主要把url后面的值变成可读
func processSurgeRule(line, hostname string) string {
	log.Println("Processing Surge rule: ", line)
	// Extract the url value
	urlValueStartIdx := strings.Index(line, "url=") + 4
	urlValueEndIdx := strings.Index(line[urlValueStartIdx:], ",")
	if urlValueEndIdx == -1 {
		urlValueEndIdx = len(line)
	} else {
		urlValueEndIdx += urlValueStartIdx
	}

	urlValue := line[urlValueStartIdx:urlValueEndIdx]

	log.Println("urlValue: ", urlValue)

	// Base64 decode the value
	decodedValue, err := base64.RawStdEncoding.DecodeString(urlValue)
	if err != nil {
		log.Println("Error decoding base64 value: ", err)
		// Handle error or just return the original line
		return line
	}

	// Extract the last part of the decoded value
	parts := strings.Split(string(decodedValue), "/")
	lastPart := parts[len(parts)-1]

	// Construct the new line
	typePart := strings.TrimSuffix(lastPart, ".list")
	newLine := "RULE-SET,https://" + hostname + "/rule?type=" + typePart

	// Append other parts from the original line
	remainingParts := strings.SplitN(line[urlValueEndIdx:], ",", 3)
	if len(remainingParts) >= 3 {
		newLine += "," + remainingParts[1] + "," + remainingParts[2]
	}

	return newLine
}
