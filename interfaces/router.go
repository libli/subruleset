package interfaces

import (
	"fmt"
	"net/http"

	"subruleset/config"
	"subruleset/interfaces/handler"
)

// SetupRouter 设置路由
func SetupRouter(config config.Config) *http.ServeMux {
	mux := http.NewServeMux()

	// 根据配置设置路由处理程序
	for key, _ := range config.Urls {
		subHandler := handler.NewSubHandler(key)
		mux.HandleFunc(fmt.Sprintf("/%s", key), subHandler.Handle)
	}

	// 为 /rule 路径添加新的处理程序
	ruleHandler := handler.NewRuleHandler()
	mux.HandleFunc("/rule", ruleHandler.Handle)

	return mux
}
