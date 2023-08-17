package config

import (
	"fmt"
	"os"
	"sync/atomic"

	"subruleset/tlog"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

// Config represents the configuration for the application
type Config struct {
	Token   string            `yaml:"token"`
	BaseURL string            `yaml:"baseURL"`
	Urls    map[string]string `yaml:"urls"`
}

var cfg atomic.Value

// Get provides a safe way to retrieve the current configuration
// 不要返回指针，如*Config，避免被修改
func Get() Config {
	return cfg.Load().(Config)
}

// Watch watches the given file for changes and updates the configuration
func Watch(filename string) error {
	tlog.Info.Printf("reading configuration from file: %s", filename)

	// Initialize the configuration for the first time
	err := readConfig(filename)
	if err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}

	// Watch for changes
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	// 无需 defer watcher.Close()，因为它应该一直运行，直到程序结束

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					tlog.Info.Println("modified file:", event.Name)
					readConfig(filename)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				tlog.Error.Println("error:", err)
			}
		}
	}()

	return watcher.Add(filename)
}

func readConfig(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var conf Config
	err = yaml.Unmarshal(data, &conf)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	cfg.Store(conf)
	return nil
}
