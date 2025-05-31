package minio

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
	"time"

	"github.com/gogodjzhu/gogobox/pkg/cmdutil"
	"github.com/minio/minio-go/v6"
	"github.com/nfnt/resize"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

type UploadOptions struct {
	Config       *MinIOConfig
	AutoResize   bool
	MaxSize      int64
	ImageQuality int
	PrintURLs   bool
}

func NewCmdMinIOUpload(f *cmdutil.Factory) *cobra.Command {
	opts := &UploadOptions{
		Config:       NewDefaultConfig(),
		AutoResize:   true,
		MaxSize:      1024 * 1024, // 1MB
		ImageQuality: 60,
		PrintURLs:   true,
	}

	cmd := &cobra.Command{
		Use:   "upload [flags] <file1> [file2] ...",
		Short: "Upload files to MinIO",
		Long: `Upload files to a MinIO server with automatic image optimization.

This command uploads files to MinIO and can automatically resize large images
to reduce file size. Supported image formats: PNG, JPG, JPEG.

The command will:
- Validate all required configuration parameters
- Process and optimize images if they exceed the size limit
- Upload files to the specified bucket
- Return public URLs for uploaded files`,
		Example: `  # Upload files with basic configuration
  gogobox minio upload -e localhost:9000 -a mykey -s mysecret -b mybucket image.jpg`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpload(opts, args)
		},
	}

	// MinIO connection flags
	cmd.Flags().StringVarP(&opts.Config.Endpoint, "endpoint", "e", "", "MinIO server endpoint (required)")
	cmd.Flags().StringVarP(&opts.Config.AccessKeyID, "access-key", "a", "", "MinIO access key ID (required)")
	cmd.Flags().StringVarP(&opts.Config.SecretAccessKey, "secret-key", "s", "", "MinIO secret access key (required)")
	cmd.Flags().StringVarP(&opts.Config.BucketName, "bucket", "b", "", "MinIO bucket name (required)")
	cmd.Flags().BoolVar(&opts.Config.UseSSL, "ssl", false, "Use SSL/TLS connection")
	cmd.Flags().StringVar(&opts.Config.Region, "region", "us-east-1", "MinIO region")

	// Upload options flags
	cmd.Flags().BoolVar(&opts.AutoResize, "resize", true, "Automatically resize large images")
	cmd.Flags().BoolVar(&opts.AutoResize, "no-resize", false, "Disable automatic image resizing")
	cmd.Flags().Int64Var(&opts.MaxSize, "max-size", 1024*1024, "Maximum file size before resizing (in bytes)")
	cmd.Flags().IntVar(&opts.ImageQuality, "quality", 60, "JPEG quality for resized images (1-100)")
	cmd.Flags().BoolVar(&opts.PrintURLs, "print-urls", true, "Print public URLs for uploaded files")

	// Mark required flags
	cmd.MarkFlagRequired("endpoint")
	cmd.MarkFlagRequired("access-key")
	cmd.MarkFlagRequired("secret-key")
	cmd.MarkFlagRequired("bucket")

	return cmd
}

func runUpload(opts *UploadOptions, filenames []string) error {
	// Validate configuration
	if err := opts.Config.Validate(); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	// Process files (resize if needed)
	processedFiles, err := processFiles(filenames, opts)
	if err != nil {
		return fmt.Errorf("file processing error: %w", err)
	}

	// Upload files
	urls, err := uploadFiles(processedFiles, opts)
	if err != nil {
		return fmt.Errorf("upload error: %w", err)
	}

	// Display results
	if opts.PrintURLs {
		fmt.Println("Upload Success:")
		for _, url := range urls {
			fmt.Printf("%s\n", url)
		}
	} else {
		fmt.Printf("Uploaded %d files successfully\n", len(urls))
	}

	return nil
}

func processFiles(filenames []string, opts *UploadOptions) ([]string, error) {
	if !opts.AutoResize {
		return filenames, nil
	}

	processedFiles := make([]string, 0, len(filenames))

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
		}

		stat, err := file.Stat()
		if err != nil {
			file.Close()
			return nil, fmt.Errorf("failed to get file stats for %s: %w", filename, err)
		}

		// If file is small enough or not an image, use original
		if stat.Size() <= opts.MaxSize || !isImage(filename) {
			file.Close()
			processedFiles = append(processedFiles, filename)
			continue
		}

		// Process large image
		processedFile, err := resizeImage(file, filename, opts)
		file.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to resize image %s: %w", filename, err)
		}

		processedFiles = append(processedFiles, processedFile)
	}

	return processedFiles, nil
}

func isImage(filename string) bool {
	lower := strings.ToLower(filename)
	return strings.HasSuffix(lower, ".png") ||
		strings.HasSuffix(lower, ".jpg") ||
		strings.HasSuffix(lower, ".jpeg")
}

func resizeImage(file *os.File, filename string, opts *UploadOptions) (string, error) {
	// Reset file pointer
	file.Seek(0, 0)

	var img image.Image
	var err error

	switch {
	case strings.HasSuffix(strings.ToLower(filename), ".png"):
		img, err = png.Decode(file)
	case strings.HasSuffix(strings.ToLower(filename), ".jpg") || strings.HasSuffix(strings.ToLower(filename), ".jpeg"):
		img, err = jpeg.Decode(file)
	default:
		return "", errors.New("unsupported image format: " + filename)
	}

	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize image
	resizedImg := resize.Resize(1080, 0, img, resize.Lanczos3)

	// Create temporary file
	tempFilename := fmt.Sprintf("/tmp/%s.jpeg", uuid.NewV1().String())
	tempFile, err := os.Create(tempFilename)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	// Encode as JPEG
	if err := jpeg.Encode(tempFile, resizedImg, &jpeg.Options{Quality: opts.ImageQuality}); err != nil {
		return "", fmt.Errorf("failed to encode resized image: %w", err)
	}

	return tempFilename, nil
}

func uploadFiles(filenames []string, opts *UploadOptions) ([]string, error) {
	// Initialize MinIO client
	minioClient, err := minio.New(opts.Config.Endpoint, opts.Config.AccessKeyID, opts.Config.SecretAccessKey, opts.Config.UseSSL)
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	var urls []string

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
		}

		fileStat, err := file.Stat()
		if err != nil {
			file.Close()
			return nil, fmt.Errorf("failed to get file stats for %s: %w", filename, err)
		}

		// Generate object name
		objectName := generateObjectName(filename)

		// Determine content type
		contentType := getContentType(filename)

		// Upload file
		_, err = minioClient.PutObject(
			opts.Config.BucketName,
			objectName,
			file,
			fileStat.Size(),
			minio.PutObjectOptions{ContentType: contentType},
		)
		file.Close()

		if err != nil {
			return nil, fmt.Errorf("failed to upload file %s: %w", filename, err)
		}

		// Generate public URL if requested
		if opts.PrintURLs {
			url := opts.Config.GetObjectURL(objectName)
			urls = append(urls, url)
		} else {
			urls = append(urls, objectName)
		}

		// Clean up temporary files
		if strings.HasPrefix(filename, "/tmp/") && strings.Contains(filename, uuid.NewV1().String()[:8]) {
			os.Remove(filename)
		}
	}

	return urls, nil
}

func generateObjectName(filename string) string {
	// Extract file extension
	parts := strings.Split(filename, ".")
	suffix := "bin"
	if len(parts) > 1 {
		suffix = parts[len(parts)-1]
	}

	// Generate unique object name with timestamp and UUID
	timestamp := time.Now().Format("20060102_150405")
	uuidStr := uuid.NewV4().String()
	return fmt.Sprintf("%s_%s.%s", timestamp, uuidStr, suffix)
}

func getContentType(filename string) string {
	lower := strings.ToLower(filename)
	switch {
	case strings.HasSuffix(lower, ".png"):
		return "image/png"
	case strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg"):
		return "image/jpeg"
	case strings.HasSuffix(lower, ".gif"):
		return "image/gif"
	case strings.HasSuffix(lower, ".pdf"):
		return "application/pdf"
	case strings.HasSuffix(lower, ".txt"):
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}
