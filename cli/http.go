package cli

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// getIgnoreData constructs a url from 'url' and 'targets', hits that URL
// and returns the response data and any error
func getIgnoreData(url string, targets []string) ([]byte, error) {
	targetString := strings.Join(targets, ",")
	fullURL := strings.Join([]string{url, targetString}, "/")

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("HTTP error: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read HTTP response: %w", err)
	}

	return data, nil
}
