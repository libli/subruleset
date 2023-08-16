package handler

import (
	"fmt"
	"net/http"
)

// copyHeaderIfExists 复制header
func copyHeaderIfExists(source, target http.Header, headerName string) {
	if values, ok := source[headerName]; ok {
		for _, value := range values {
			target.Add(headerName, value)
		}
	}
}

// getCurrentURL 获取当前URL
func getCurrentURL(r *http.Request) string {
	schema := "http" // 默认值
	if forwardedProto := r.Header.Get("X-Forwarded-Proto"); forwardedProto != "" {
		schema = forwardedProto
	}
	host := r.Host
	hostname := fmt.Sprintf("%s://%s", schema, host)
	return fmt.Sprintf("%s%s", hostname, r.URL)
}
