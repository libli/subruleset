package logic

import (
	"fmt"
	"io"
	"net/http"
)

// getHeadersAndContentFromURL 从URL获取内容
func getHeadersAndContentFromURL(url string) (http.Header, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, "", fmt.Errorf("error fetching data from URL: %s, err: %v", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("received unexpected status: %s", resp.Status)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("error reading response body: %v", err)
	}

	return resp.Header, string(content), nil
}
