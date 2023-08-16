package handler

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"

	"subruleset/tlog"
)

type RuleHandler struct {
}

// NewRuleHandler 创建一个新的NewRuleHandler
func NewRuleHandler() *RuleHandler {
	return &RuleHandler{}
}

// Handle 处理 /rule 路径的请求
func (h *RuleHandler) Handle(w http.ResponseWriter, r *http.Request) {
	tlog.Info.Printf("Handling request for URL: %s", r.URL)
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
		tlog.Error.Printf("Error fetching data from URL: %s, err: %v", newURL, err)
		http.Error(w, "Failed to get content from the URL", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	content, err := io.ReadAll(response.Body)
	if err != nil {
		tlog.Error.Printf("Error reading response body: %v", err)
		http.Error(w, "Failed to read content from the response", http.StatusInternalServerError)
		return
	}

	// 将内容返回给用户
	w.Write(content)
}
