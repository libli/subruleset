package logic

import (
	"encoding/base64"
	"fmt"

	"subruleset/config"
)

const (
	path          = "/getruleset"
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
	config := config.Get()
	rulePath := fmt.Sprintf(rulesPathTmpl, ruleType)
	encodedPath := base64.RawStdEncoding.EncodeToString([]byte(rulePath))
	newURL := fmt.Sprintf("%s%s?type=1&url=%s", config.BaseURL, path, encodedPath)
	_, content, err := getHeadersAndContentFromURL(newURL)
	return content, err
}
