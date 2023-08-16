package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"subruleset/tlog"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Token string            `yaml:"token"`
	Urls  map[string]string `yaml:"urls"`
}

func main() {
	config, err := readConfig("config.yaml")
	if err != nil {
		tlog.Error.Fatalf("Failed to read config: %v", err)
	}

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
	http.ListenAndServe(":8080", mux)
}

// readConfig 读取配置文件
func readConfig(filename string) (Config, error) {
	tlog.Info.Printf("Reading configuration from file: %s", filename)

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

	tlog.Info.Println("Successfully read configuration")
	return config, nil
}
