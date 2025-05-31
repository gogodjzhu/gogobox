package minio

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gogodjzhu/gogobox/internal/util"
	"github.com/gogodjzhu/gogobox/pkg/cmdutil"
	"github.com/minio/minio-go/v6"
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
)

type UploadOptions struct {
	Config     *MinIOConfig
	AutoResize bool
	MaxSize    int64
	PrintURLs  bool
}

func NewCmdMinIOUpload(f *cmdutil.Factory) *cobra.Command {
	opts := &UploadOptions{
		Config:     NewDefaultConfig(),
		AutoResize: true,
		MaxSize:    512 * 1024, // 512KB default max size for images
		PrintURLs:  true,
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

	// Upload options flags
	cmd.Flags().BoolVar(&opts.AutoResize, "resize", true, "Automatically resize large images")
	cmd.Flags().Int64Var(&opts.MaxSize, "max-size", 512*1024, "Maximum file size in bytes after resize")
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
		processedFile, err := util.CompressImage(filename, opts.MaxSize, "jpeg")
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

func uploadFiles(filenames []string, opts *UploadOptions) ([]string, error) {
	// Initialize MinIO client
	minioClient, err := minio.New(opts.Config.Endpoint, opts.Config.AccessKeyID, opts.Config.SecretAccessKey, opts.Config.UseSSL)
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Check if bucket exists
	exists, err := minioClient.BucketExists(opts.Config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("bucket '%s' does not exist", opts.Config.BucketName)
	}

	var urls []string
	var uploadedObjects []string // Keep track of successfully uploaded objects for potential rollback

	for _, filename := range filenames {
		file, err := os.Open(filename)
		if err != nil {
			// Clean up any already uploaded files on error
			cleanupUploadedFiles(minioClient, opts.Config.BucketName, uploadedObjects)
			return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
		}

		fileStat, err := file.Stat()
		if err != nil {
			file.Close()
			// Clean up any already uploaded files on error
			cleanupUploadedFiles(minioClient, opts.Config.BucketName, uploadedObjects)
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

		// Close file immediately after upload
		file.Close()

		if err != nil {
			// Clean up any already uploaded files on error
			cleanupUploadedFiles(minioClient, opts.Config.BucketName, uploadedObjects)
			return nil, fmt.Errorf("failed to upload file %s: %w", filename, err)
		}

		// Track successfully uploaded object
		uploadedObjects = append(uploadedObjects, objectName)

		// Generate public URL if requested
		if opts.PrintURLs {
			url := opts.Config.GetObjectURL(objectName)
			urls = append(urls, url)
		} else {
			urls = append(urls, objectName)
		}

		// Clean up temporary files
		if strings.HasPrefix(filename, "/tmp/") {
			os.Remove(filename)
		}
	}

	return urls, nil
}

// cleanupUploadedFiles removes objects that were successfully uploaded before an error occurred
func cleanupUploadedFiles(client *minio.Client, bucketName string, objectNames []string) {
	for _, objectName := range objectNames {
		// Best effort cleanup - don't propagate errors from cleanup
		client.RemoveObject(bucketName, objectName)
	}
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
