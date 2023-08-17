package logic

import (
	"encoding/base64"
	"fmt"
)

const (
	baseURL       = "https://sub.zaptiah.com/getruleset"
	rulesPathTmpl = "rules/myclash/Merge/%s.list"
)

type RuleLogic struct {
}

// NewRuleLogic 创建一个新的RuleLogic
func NewRuleLogic() *RuleLogic {
	return &RuleLogic{}
}

// FetchRuleSet 获取规则集
func (l *RuleLogic) FetchRuleSet(ruleType string) (string, error) {
	path := fmt.Sprintf(rulesPathTmpl, ruleType)
	encodedPath := base64.RawStdEncoding.EncodeToString([]byte(path))
	newURL := fmt.Sprintf("%s?type=1&url=%s", baseURL, encodedPath)
	_, content, err := getHeadersAndContentFromURL(newURL)
	return content, err
}
