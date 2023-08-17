package handler

import (
	"net/http"
	"strings"

	"subruleset/logic"
	"subruleset/tlog"
)

type SubHandler struct {
	logic *logic.SubLogic
}

// NewSubHandler 创建一个新的SubHandler
func NewSubHandler(subKey string) *SubHandler {
	return &SubHandler{
		logic: logic.NewSubLogic(subKey),
	}
}

// Handle 处理请求
func (h *SubHandler) Handle(w http.ResponseWriter, r *http.Request) {
	tlog.Info.Printf("Handling request for URL: %s", r.URL)

	token := r.URL.Query().Get("token")
	if !h.logic.ValidateToken(token) {
		tlog.Warn.Printf("Forbidden request made with token: %s", token)
		http.Error(w, "无权限", http.StatusForbidden)
		return
	}
	headers, body, err := h.logic.FetchSubscriptions()
	if err != nil {
		tlog.Error.Printf("Error fetching data from URL err: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 设置 "access-control-allow-origin" 为 "*"；设置订阅信息
	w.Header().Set("Access-Control-Allow-Origin", "*")
	copyHeaderIfExists(headers, w.Header(), "Subscription-Userinfo")

	currentURL := getCurrentURL(r)

	var newBody string
	// Check if URL path contains "-surge" substring
	if strings.Contains(r.URL.Path, "-surge") {
		newBody = h.logic.Surge(body, currentURL)
	} else {
		newBody = h.logic.Clash(body)
	}

	w.Write([]byte(newBody))
}
