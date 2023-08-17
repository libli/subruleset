package handler

import (
	"net/http"

	"subruleset/logic"
	"subruleset/tlog"
)

type RuleHandler struct {
	logic *logic.RuleLogic
}

// NewRuleHandler 创建一个新的NewRuleHandler
func NewRuleHandler() *RuleHandler {
	return &RuleHandler{
		logic: logic.NewRuleLogic(),
	}
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

	content, err := h.logic.FetchRuleSet(ruleType)
	if err != nil {
		tlog.Error.Printf("Error fetching data from type: %s, err: %v", ruleType, err)
		http.Error(w, "Failed to get content from the URL", http.StatusInternalServerError)
		return
	}

	// 将内容返回给用户
	if _, err := w.Write([]byte(content)); err != nil {
		tlog.Error.Printf("Error writing response body: %v", err)
	}
}
