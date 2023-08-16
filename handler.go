package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"subruleset/tlog"
)

// handler 处理请求
func handler(url string, confToken string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		tlog.Info.Printf("Handling request for URL: %s", r.URL)

		token := r.URL.Query().Get("token")
		if token != confToken {
			tlog.Warn.Printf("Forbidden request made with token: %s", token)
			http.Error(w, "无权限", http.StatusForbidden)
			return
		}
		resp, err := http.Get(url)
		if err != nil {
			tlog.Error.Printf("Error fetching data from URL: %s, err: %v", url, err)
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
			tlog.Error.Printf("Error reading response body: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		schema := "http" // 默认值
		if forwardedProto := r.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
			schema = forwardedProto
		}

		host := r.Host
		hostname := fmt.Sprintf("%s://%s", schema, host)
		fullURL := fmt.Sprintf("%s%s", hostname, r.URL)
		newBody := processRules(string(body), hostname, fullURL)

		w.Write([]byte(newBody))
	}
}

// ruleHandler 处理 /rule 路径的请求
func ruleHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
	ruleType := values.Get("type")
	if ruleType == "" {
		http.Error(w, "type is required", http.StatusBadRequest)
		return
	}

	// 拼接URL路径
	path := fmt.Sprintf("rules/myclash/Merge/%s.list", ruleType)
	// 进行Base64编码
	encodedPath := base64.RawStdEncoding.EncodeToString([]byte(path))

	// 使用Base64编码拼接新的URL
	newURL := fmt.Sprintf("https://sub.zaptiah.com/getruleset?type=1&url=%s", encodedPath)

	// 访问新的URL并获取其内容
	response, err := http.Get(newURL)
	if err != nil {
		http.Error(w, "Failed to get content from the URL", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Failed to read content from the response", http.StatusInternalServerError)
		return
	}

	// 将内容返回给用户
	w.Write(content)
}
