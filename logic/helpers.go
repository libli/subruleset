package logic

import (
	"fmt"
	"io"
	"net/http"

	"subruleset/tlog"
)

// getHeadersAndContentFromURL 从URL获取内容
func getHeadersAndContentFromURL(url string) (http.Header, string, error) {
	tlog.Info.Printf("Fetching data from URL: %s", url)
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

	tlog.Info.Printf("Successfully fetched data from URL: %s", url)
	return resp.Header, string(content), nil
}
