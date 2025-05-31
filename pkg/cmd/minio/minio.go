package minio

import (
	"errors"
	"fmt"

	"github.com/gogodjzhu/gogobox/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdMinIO(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "minio",
		Short: "MinIO operations",
		Long: `MinIO operations for file upload and management.

This command provides various MinIO operations including file upload,
download, and bucket management. Use subcommands to perform specific operations.`,
		Run: func(cmd *cobra.Command, args []string) {
			// Show help when no subcommand is provided
			cmd.Help()
		},
	}

	// Add subcommands
	cmd.AddCommand(NewCmdMinIOUpload(f))

	return cmd
}

// MinIOConfig represents the configuration for MinIO operations
type MinIOConfig struct {
	// Endpoint is the MinIO server endpoint (e.g., "localhost:9000" or "s3.amazonaws.com")
	Endpoint string `json:"endpoint" yaml:"endpoint"`

	// AccessKeyID is the access key for MinIO authentication
	AccessKeyID string `json:"accessKeyID" yaml:"accessKeyID"`

	// SecretAccessKey is the secret key for MinIO authentication
	SecretAccessKey string `json:"secretAccessKey" yaml:"secretAccessKey"`

	// BucketName is the name of the bucket to upload files to
	BucketName string `json:"bucketName" yaml:"bucketName"`

	// UseSSL indicates whether to use HTTPS for connection
	UseSSL bool `json:"useSSL" yaml:"useSSL"`
}

// Validate checks if the MinIO configuration is valid
func (c *MinIOConfig) Validate() error {
	if c.Endpoint == "" {
		return errors.New("endpoint must not be empty")
	}
	if c.AccessKeyID == "" {
		return errors.New("accessKeyID must not be empty")
	}
	if c.SecretAccessKey == "" {
		return errors.New("secretAccessKey must not be empty")
	}
	if c.BucketName == "" {
		return errors.New("bucketName must not be empty")
	}
	return nil
}

// GetObjectURL returns the public URL for an object in the bucket
func (c *MinIOConfig) GetObjectURL(objectName string) string {
	protocol := "http"
	if c.UseSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, c.Endpoint, c.BucketName, objectName)
}

// NewDefaultConfig returns a new MinIOConfig with default values
func NewDefaultConfig() *MinIOConfig {
	return &MinIOConfig{
		UseSSL: false,
	}
}
