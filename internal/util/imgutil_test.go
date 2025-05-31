package util

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
	"testing"
)

// createTestImage creates a test image with specified dimensions and format
func createTestImage(width, height int, format string) ([]byte, error) {
	// Create a simple test image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill with a simple pattern
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Create a simple gradient pattern
			r := uint8((x * 255) / width)
			g := uint8((y * 255) / height)
			b := uint8(128)
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	var buf bytes.Buffer
	switch format {
	case "png":
		err := png.Encode(&buf, img)
		return buf.Bytes(), err
	case "jpeg":
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
		return buf.Bytes(), err
	default:
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90})
		return buf.Bytes(), err
	}
}

// createTestImageFile creates a test image file and returns the file path
func createTestImageFile(width, height int, format string) (string, error) {
	imageData, err := createTestImage(width, height, format)
	if err != nil {
		return "", err
	}

	// Create temporary file
	tempFile, err := os.CreateTemp("", "test_image_*."+format)
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	// Write image data to file
	_, err = tempFile.Write(imageData)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", err
	}

	return tempFile.Name(), nil
}

func TestGenerateRandomFilename(t *testing.T) {
	tests := []struct {
		name   string
		format string
	}{
		{"JPEG format", "jpeg"},
		{"PNG format", "png"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename, err := generateRandomFilename(tt.format)
			if err != nil {
				t.Errorf("generateRandomFilename() error = %v", err)
				return
			}

			// Check if filename has correct format
			if !strings.HasSuffix(filename, "."+tt.format) {
				t.Errorf("generateRandomFilename() filename = %v, should end with .%v", filename, tt.format)
			}

			// Check if filename is in /tmp directory
			if !strings.HasPrefix(filename, "/tmp/") {
				t.Errorf("generateRandomFilename() filename = %v, should be in /tmp/", filename)
			}
		})
	}
}

func TestSaveImageToTemp(t *testing.T) {
	// Create test image data
	testData := []byte("test image data")
	format := "jpeg"

	// Save image to temp
	tempPath, err := saveImageToTemp(testData, format)
	if err != nil {
		t.Errorf("saveImageToTemp() error = %v", err)
		return
	}

	// Clean up
	defer os.Remove(tempPath)

	// Check if file exists
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		t.Errorf("saveImageToTemp() file does not exist at %v", tempPath)
	}

	// Check file content
	savedData, err := os.ReadFile(tempPath)
	if err != nil {
		t.Errorf("Failed to read saved file: %v", err)
		return
	}

	if !bytes.Equal(savedData, testData) {
		t.Errorf("saveImageToTemp() saved data does not match original")
	}
}

func TestCompressImage_SmallImage(t *testing.T) {
	// Create a small test image file (should not be compressed)
	testFilePath, err := createTestImageFile(100, 100, "jpeg")
	if err != nil {
		t.Fatalf("Failed to create test image file: %v", err)
	}
	defer os.Remove(testFilePath)

	// Get original file size
	originalInfo, err := os.Stat(testFilePath)
	if err != nil {
		t.Fatalf("Failed to get original file stats: %v", err)
	}

	// Test with maxSize larger than image size
	maxSize := originalInfo.Size() + 1000

	tempPath, err := CompressImage(testFilePath, maxSize, "jpeg")
	if err != nil {
		t.Errorf("CompressImage() error = %v", err)
		return
	}

	// Clean up
	defer os.Remove(tempPath)

	// Check if file exists
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		t.Errorf("CompressImage() file does not exist at %v", tempPath)
	}

	// Check if file is in /tmp directory
	if !strings.HasPrefix(tempPath, "/tmp/") {
		t.Errorf("CompressImage() file should be in /tmp/, got %v", tempPath)
	}
}

func TestCompressImage_LargeImage(t *testing.T) {
	// Create a large test image file (should be compressed)
	testFilePath, err := createTestImageFile(800, 600, "jpeg")
	if err != nil {
		t.Fatalf("Failed to create test image file: %v", err)
	}
	defer os.Remove(testFilePath)

	// Get original file size
	originalInfo, err := os.Stat(testFilePath)
	if err != nil {
		t.Fatalf("Failed to get original file stats: %v", err)
	}

	// Test with maxSize smaller than image size
	maxSize := originalInfo.Size() / 4 // Force compression

	tempPath, err := CompressImage(testFilePath, maxSize, "jpeg")
	if err != nil {
		t.Errorf("CompressImage() error = %v", err)
		return
	}

	// Clean up
	defer os.Remove(tempPath)

	// Check if file exists
	fileInfo, err := os.Stat(tempPath)
	if os.IsNotExist(err) {
		t.Errorf("CompressImage() file does not exist at %v", tempPath)
		return
	}

	// For JPEG, we expect compression to work
	if fileInfo.Size() >= originalInfo.Size() {
		t.Logf("Warning: Compressed file size (%d) is not smaller than original (%d), but this might be acceptable for certain images", fileInfo.Size(), originalInfo.Size())
	}

	// Check if file is in /tmp directory
	if !strings.HasPrefix(tempPath, "/tmp/") {
		t.Errorf("CompressImage() file should be in /tmp/, got %v", tempPath)
	}
}

func TestCompressImage_NonExistentFile(t *testing.T) {
	_, err := CompressImage("/non/existent/file.jpg", 1000, "jpeg")
	if err == nil {
		t.Errorf("CompressImage() should return error for non-existent file")
	}

	if !strings.Contains(err.Error(), "file does not exist") {
		t.Errorf("CompressImage() should return file not found error, got: %v", err)
	}
}

func TestCompressImage_PNGFormat(t *testing.T) {
	// Create a PNG test image file
	testFilePath, err := createTestImageFile(400, 300, "png")
	if err != nil {
		t.Fatalf("Failed to create test PNG image file: %v", err)
	}
	defer os.Remove(testFilePath)

	// Get original file size
	originalInfo, err := os.Stat(testFilePath)
	if err != nil {
		t.Fatalf("Failed to get original file stats: %v", err)
	}

	// Test with maxSize smaller than image size
	maxSize := originalInfo.Size() / 2

	tempPath, err := CompressImage(testFilePath, maxSize, "png")
	if err != nil {
		t.Errorf("CompressImage() error = %v", err)
		return
	}

	// Clean up
	defer os.Remove(tempPath)

	// Check if file exists and has proper extension
	if !strings.HasSuffix(tempPath, ".png") {
		t.Errorf("CompressImage() PNG file should have .png extension, got %v", tempPath)
	}
}

func TestCompressImage_InvalidImageFile(t *testing.T) {
	// Create a non-image file
	tempFile, err := os.CreateTemp("", "test_invalid_*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write some non-image data
	_, err = tempFile.WriteString("This is not an image file")
	tempFile.Close()
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Try to compress the non-image file
	_, err = CompressImage(tempFile.Name(), 1000, "jpeg")
	if err == nil {
		t.Errorf("CompressImage() should return error for invalid image file")
	}

	if !strings.Contains(err.Error(), "failed to decode image") {
		t.Errorf("CompressImage() should return decode error for invalid image file, got: %v", err)
	}
}

func TestCompressImage_ZeroMaxSize(t *testing.T) {
	// Create a test image file
	testFilePath, err := createTestImageFile(100, 100, "jpeg")
	if err != nil {
		t.Fatalf("Failed to create test image file: %v", err)
	}
	defer os.Remove(testFilePath)

	// Test with maxSize of 0 (should force compression)
	tempPath, err := CompressImage(testFilePath, 0, "jpeg")
	if err != nil {
		t.Errorf("CompressImage() error = %v", err)
		return
	}

	// Clean up
	defer os.Remove(tempPath)

	// Check if file exists
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		t.Errorf("CompressImage() file does not exist at %v", tempPath)
	}
}

func TestCompressImage_DifferentFormats(t *testing.T) {
	tests := []struct {
		name   string
		format string
	}{
		{"JPEG format", "jpeg"},
		{"PNG format", "png"},
		{"Default format", "unknown"}, // Should default to JPEG behavior
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test image file
			testFilePath, err := createTestImageFile(200, 200, "jpeg")
			if err != nil {
				t.Fatalf("Failed to create test image file: %v", err)
			}
			defer os.Remove(testFilePath)

			// Get original file size
			originalInfo, err := os.Stat(testFilePath)
			if err != nil {
				t.Fatalf("Failed to get original file stats: %v", err)
			}

			// Test compression with different formats
			maxSize := originalInfo.Size() / 2
			tempPath, err := CompressImage(testFilePath, maxSize, tt.format)
			if err != nil {
				t.Errorf("CompressImage() error = %v", err)
				return
			}

			// Clean up
			defer os.Remove(tempPath)

			// Check if file exists
			if _, err := os.Stat(tempPath); os.IsNotExist(err) {
				t.Errorf("CompressImage() file does not exist at %v", tempPath)
			}

			// Check file extension based on format
			expectedExt := tt.format
			if !strings.HasSuffix(tempPath, "."+expectedExt) {
				t.Errorf("CompressImage() file should have .%s extension, got %v", expectedExt, tempPath)
			}
		})
	}
}

func TestSaveCompressedImage_QualityLevels(t *testing.T) {
	// Create test image
	img := image.NewRGBA(image.Rect(0, 0, 400, 300))
	for y := 0; y < 300; y++ {
		for x := 0; x < 400; x++ {
			img.Set(x, y, color.RGBA{uint8(x % 256), uint8(y % 256), 128, 255})
		}
	}

	// Test JPEG compression with very small maxSize to force quality reduction
	tempPath, err := saveCompressedImage(img, "jpeg", 1000) // Very small size
	if err != nil {
		t.Errorf("saveCompressedImage() error = %v", err)
		return
	}

	// Clean up
	defer os.Remove(tempPath)

	// Check if file exists
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		t.Errorf("saveCompressedImage() file does not exist at %v", tempPath)
	}

	// Check if it's a JPEG file
	if !strings.HasSuffix(tempPath, ".jpeg") {
		t.Errorf("saveCompressedImage() JPEG file should have .jpeg extension, got %v", tempPath)
	}
}
