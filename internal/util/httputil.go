package util

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
)

func SendGet(url string, header map[string]string, wrap func(response *http.Response) (interface{}, error)) (interface{}, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	host, err := hostname(url)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Host", host)
	for k, v := range header {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return wrap(resp)
}

func SendPost(url string, header map[string]string, body []byte, wrap func(response *http.Response) (interface{}, error)) (interface{}, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Length", strconv.Itoa(len(body)))
	host, err := hostname(url)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Host", host)
	for k, v := range header {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return wrap(resp)
}

func hostname(u string) (string, error) {
	e, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	return e.Hostname(), nil
}

// DownloadToTempFile downloads a file from the given URL and saves it to a temporary file in /tmp/.
// It returns the path to the temporary file.
func DownloadToTempFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download %s: status %s", url, resp.Status)
	}

	ext := filepath.Ext(url)
	if ext == "" {
		ext = ".bin"
	}
	tmpFile, err := os.CreateTemp("/tmp", "gogobox_download_*"+ext)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("failed to write to temp file: %w", err)
	}

	return tmpFile.Name(), nil
}
