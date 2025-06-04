package util

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestDownloadToTempFile(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello world"))
	}))
	defer ts.Close()

	url := ts.URL + "/testfile.txt"
	path, err := DownloadToTempFile(url)
	if err != nil {
		t.Fatalf("DownloadToTempFile failed: %v", err)
	}
	defer os.Remove(path)

	// Check file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Downloaded file does not exist: %v", path)
	}

	// Check file content
	content, err := ioutil.ReadFile(path)
	if err != nil {
		t.Errorf("Failed to read downloaded file: %v", err)
	}
	if string(content) != "hello world" {
		t.Errorf("Unexpected file content: %s", string(content))
	}
}
