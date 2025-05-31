package util

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

// CompressImage reads a local image file and compresses it if it exceeds maxSize.
// The compressed image is saved to /tmp/{random-name}.{format} while maintaining aspect ratio.
// If the original image is smaller than maxSize, it's saved without compression.
func CompressImage(filepath string, maxSize int64, format string) (string, error) {
	// Check if file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return "", fmt.Errorf("file does not exist: %s", filepath)
	}

	// Read the image file
	imgData, err := os.ReadFile(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to read image file %s: %w", filepath, err)
	}

	// Validate that the file is actually an image by trying to decode it
	_, _, err = image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return "", fmt.Errorf("failed to decode image file %s: %w", filepath, err)
	}

	// Check if the original image is already within the size limit
	if int64(len(imgData)) <= maxSize {
		// Save the original image without compression
		return saveImageToTemp(imgData, format)
	}

	// Compress the image
	return compressAndSaveImage(imgData, maxSize, format)
}

// generateRandomFilename generates a random filename for temporary files
func generateRandomFilename(format string) (string, error) {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	filename := fmt.Sprintf("%x.%s", randomBytes, format)
	return filepath.Join("/tmp", filename), nil
}

// saveImageToTemp saves image data to a temporary file
func saveImageToTemp(data []byte, format string) (string, error) {
	tempPath, err := generateRandomFilename(format)
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(tempPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write image to temp file: %w", err)
	}

	return tempPath, nil
}

// compressAndSaveImage compresses the image and saves it to a temporary file
func compressAndSaveImage(imgData []byte, maxSize int64, format string) (string, error) {
	// Decode the image
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Try different compression qualities to meet the size requirement
	return saveCompressedImage(img, format, maxSize)
}

// saveCompressedImage saves the resized image with appropriate compression
func saveCompressedImage(img image.Image, format string, maxSize int64) (string, error) {
	qualities := []int{90, 80, 70, 60, 50, 40, 30, 20, 10} // JPEG quality levels to try

	for _, quality := range qualities {
		// Generate temporary file path
		tempPath, err := generateRandomFilename(format)
		if err != nil {
			return "", err
		}

		// Create and encode the image
		file, err := os.Create(tempPath)
		if err != nil {
			return "", fmt.Errorf("failed to create temp file: %w", err)
		}

		var encodeErr error
		switch format {
		case "png":
			// PNG doesn't have quality settings, so we encode once
			encodeErr = png.Encode(file, img)
		case "jpeg":
			encodeErr = jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
		default:
			// Default to JPEG
			encodeErr = jpeg.Encode(file, img, &jpeg.Options{Quality: quality})
		}

		file.Close()

		if encodeErr != nil {
			os.Remove(tempPath) // Clean up on error
			return "", fmt.Errorf("failed to encode image: %w", encodeErr)
		}

		// Check the file size
		fileInfo, err := os.Stat(tempPath)
		if err != nil {
			os.Remove(tempPath) // Clean up on error
			return "", fmt.Errorf("failed to get file info: %w", err)
		}

		// If the file size is within the limit, return the path
		if fileInfo.Size() <= maxSize {
			return tempPath, nil
		}

		// If PNG and still too large, we can't compress further with quality
		if format == "png" {
			return tempPath, nil // Return the PNG as-is
		}

		// Clean up and try next quality level for JPEG
		os.Remove(tempPath)
	}

	// If all quality levels failed, return the lowest quality version
	tempPath, err := generateRandomFilename(format)
	if err != nil {
		return "", err
	}

	file, err := os.Create(tempPath)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer file.Close()

	// Use the lowest quality
	if format == "png" {
		err = png.Encode(file, img)
	} else {
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 5})
	}

	if err != nil {
		return "", fmt.Errorf("failed to encode final image: %w", err)
	}

	return tempPath, nil
}
