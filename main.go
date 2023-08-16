package main

import (
	"fmt"
	"net/http"

	"subruleset/config"
	"subruleset/tlog"
)

func main() {
	err := config.Watch("config.yaml")
	if err != nil {
		tlog.Error.Fatalf("Failed to watch config: %v", err)
	}
	config := config.Get()

	if config.Token == "" {
		tlog.Error.Fatal("Token is empty")
	}

	mux := http.NewServeMux()
	for key, url := range config.Urls {
		tlog.Info.Printf("Adding handler for path: %s, url: %s", key, url)
		mux.HandleFunc(fmt.Sprintf("/%s", key), handler(url, config.Token))
	}

	// 为 /rule 路径添加新的处理程序
	mux.HandleFunc("/rule", ruleHandler)

	tlog.Info.Println("Starting server on port 8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		tlog.Error.Fatalf("Failed to start server: %v", err)
	}
}
