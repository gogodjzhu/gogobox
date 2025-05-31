package minio

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestIsImage tests the isImage function
func TestIsImage(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{"PNG file", "test.png", true},
		{"PNG file uppercase", "test.PNG", true},
		{"JPG file", "test.jpg", true},
		{"JPG file uppercase", "test.JPG", true},
		{"JPEG file", "test.jpeg", true},
		{"JPEG file uppercase", "test.JPEG", true},
		{"Mixed case PNG", "test.Png", true},
		{"Mixed case JPG", "test.JpG", true},
		{"Text file", "test.txt", false},
		{"PDF file", "test.pdf", false},
		{"No extension", "test", false},
		{"Empty filename", "", false},
		{"Multiple dots", "test.backup.png", true},
		{"Multiple dots non-image", "test.backup.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isImage(tt.filename); got != tt.want {
				t.Errorf("isImage(%s) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

// TestGetContentType tests the getContentType function
func TestGetContentType(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{"PNG file", "test.png", "image/png"},
		{"PNG file uppercase", "test.PNG", "image/png"},
		{"JPG file", "test.jpg", "image/jpeg"},
		{"JPG file uppercase", "test.JPG", "image/jpeg"},
		{"JPEG file", "test.jpeg", "image/jpeg"},
		{"JPEG file uppercase", "test.JPEG", "image/jpeg"},
		{"GIF file", "test.gif", "image/gif"},
		{"GIF file uppercase", "test.GIF", "image/gif"},
		{"PDF file", "test.pdf", "application/pdf"},
		{"PDF file uppercase", "test.PDF", "application/pdf"},
		{"Text file", "test.txt", "text/plain"},
		{"Text file uppercase", "test.TXT", "text/plain"},
		{"Unknown extension", "test.xyz", "application/octet-stream"},
		{"No extension", "test", "application/octet-stream"},
		{"Empty filename", "", "application/octet-stream"},
		{"Multiple dots", "test.backup.png", "image/png"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getContentType(tt.filename); got != tt.want {
				t.Errorf("getContentType(%s) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

// TestGenerateObjectName tests the generateObjectName function
func TestGenerateObjectName(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{"PNG file", "test.png"},
		{"JPG file", "test.jpg"},
		{"Text file", "test.txt"},
		{"No extension", "test"},
		{"Multiple dots", "test.backup.png"},
		{"Path with directory", "/path/to/test.png"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateObjectName(tt.filename)
			
			// Check that the result is not empty
			if got == "" {
				t.Errorf("generateObjectName(%s) returned empty string", tt.filename)
			}
			
			// Check that the result contains a timestamp pattern (YYYYMMDD_HHMMSS)
			if !strings.Contains(got, "_") {
				t.Errorf("generateObjectName(%s) = %s, expected to contain timestamp", tt.filename, got)
			}
			
			// Check that the result contains a UUID-like pattern (contains hyphens)
			if !strings.Contains(got, "-") {
				t.Errorf("generateObjectName(%s) = %s, expected to contain UUID", tt.filename, got)
			}
			
			// Check that the extension is preserved correctly
			parts := strings.Split(tt.filename, ".")
			expectedSuffix := "bin"
			if len(parts) > 1 {
				expectedSuffix = parts[len(parts)-1]
			}
			
			if !strings.HasSuffix(got, "."+expectedSuffix) {
				t.Errorf("generateObjectName(%s) = %s, expected to end with .%s", tt.filename, got, expectedSuffix)
			}
		})
	}
}

// TestProcessFiles tests the processFiles function
func TestProcessFiles(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := ioutil.TempDir("", "upload_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	smallFile := filepath.Join(tmpDir, "small.txt")
	if err := ioutil.WriteFile(smallFile, []byte("small content"), 0644); err != nil {
		t.Fatalf("Failed to create small test file: %v", err)
	}

	largeFile := filepath.Join(tmpDir, "large.txt")
	largeContent := strings.Repeat("large content", 1000) // Create a file larger than 1KB
	if err := ioutil.WriteFile(largeFile, []byte(largeContent), 0644); err != nil {
		t.Fatalf("Failed to create large test file: %v", err)
	}

	nonExistentFile := filepath.Join(tmpDir, "nonexistent.txt")

	tests := []struct {
		name          string
		filenames     []string
		opts          *UploadOptions
		expectError   bool
		expectedCount int
	}{
		{
			name:      "AutoResize disabled",
			filenames: []string{smallFile, largeFile},
			opts: &UploadOptions{
				AutoResize: false,
				MaxSize:    1024,
			},
			expectError:   false,
			expectedCount: 2,
		},
		{
			name:      "Small file under limit",
			filenames: []string{smallFile},
			opts: &UploadOptions{
				AutoResize: true,
				MaxSize:    1024 * 1024, // 1MB
			},
			expectError:   false,
			expectedCount: 1,
		},
		{
			name:      "Large non-image file",
			filenames: []string{largeFile},
			opts: &UploadOptions{
				AutoResize: true,
				MaxSize:    1024, // 1KB
			},
			expectError:   false,
			expectedCount: 1, // Non-image files are not processed
		},
		{
			name:      "Nonexistent file",
			filenames: []string{nonExistentFile},
			opts: &UploadOptions{
				AutoResize: true,
				MaxSize:    1024,
			},
			expectError:   true,
			expectedCount: 0,
		},
		{
			name:      "Mixed files",
			filenames: []string{smallFile, largeFile},
			opts: &UploadOptions{
				AutoResize: true,
				MaxSize:    1024, // 1KB
			},
			expectError:   false,
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processFiles(tt.filenames, tt.opts)
			
			if tt.expectError {
				if err == nil {
					t.Errorf("processFiles() expected error, but got none")
				}
				return
			}
			
			if err != nil {
				t.Errorf("processFiles() unexpected error: %v", err)
				return
			}
			
			if len(result) != tt.expectedCount {
				t.Errorf("processFiles() returned %d files, expected %d", len(result), tt.expectedCount)
			}
		})
	}
}

// Test edge cases
func TestIsImageEdgeCases(t *testing.T) {
	edgeCases := []struct {
		name     string
		filename string
		want     bool
	}{
		{"Filename with only dot", ".", false},
		{"Filename with multiple consecutive dots", "test..png", true},
		{"Filename ending with dot", "test.", false},
		{"Very long filename", strings.Repeat("a", 1000) + ".png", true},
		{"Filename with spaces", "test file.png", true},
		{"Filename with special chars", "test@#$.png", true},
		{"Filename with unicode", "测试.png", true},
	}

	for _, tt := range edgeCases {
		t.Run(tt.name, func(t *testing.T) {
			if got := isImage(tt.filename); got != tt.want {
				t.Errorf("isImage(%s) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

func TestGetContentTypeEdgeCases(t *testing.T) {
	edgeCases := []struct {
		name     string
		filename string
		want     string
	}{
		{"Filename with only dot", ".", "application/octet-stream"},
		{"Filename with multiple consecutive dots", "test..png", "image/png"},
		{"Filename ending with dot", "test.", "application/octet-stream"},
		{"Very long filename", strings.Repeat("a", 1000) + ".pdf", "application/pdf"},
		{"Filename with spaces", "test file.txt", "text/plain"},
		{"Filename with special chars", "test@#$.jpg", "image/jpeg"},
		{"Filename with unicode", "测试.gif", "image/gif"},
	}

	for _, tt := range edgeCases {
		t.Run(tt.name, func(t *testing.T) {
			if got := getContentType(tt.filename); got != tt.want {
				t.Errorf("getContentType(%s) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}
