package util

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3Uploader handles file uploads to AWS S3
type S3Uploader struct {
	client     *s3.Client
	bucketName string
	region     string
}

// NewS3Uploader creates a new S3Uploader instance
func NewS3Uploader(bucketName, region string) (*S3Uploader, error) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	// Create S3 client
	client := s3.NewFromConfig(cfg)

	return &S3Uploader{
		client:     client,
		bucketName: bucketName,
		region:     region,
	}, nil
}

// UploadFile uploads a file to S3 and returns the URL
func (u *S3Uploader) UploadFile(file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// Read file content
	buffer := make([]byte, fileHeader.Size)
	_, err := file.Read(buffer)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Reset file pointer
	file.Seek(0, 0)

	// Generate unique filename
	fileExt := filepath.Ext(fileHeader.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), fileExt)

	// Determine content type
	contentType := http.DetectContentType(buffer)

	// Upload to S3
	_, err = u.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(u.bucketName),
		Key:         aws.String(fileName),
		Body:        bytes.NewReader(buffer),
		ContentType: aws.String(contentType),
		ACL:         aws.String("public-read"),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	// Generate URL
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.bucketName, u.region, fileName)
	return fileURL, nil
}

// DeleteFile deletes a file from S3
func (u *S3Uploader) DeleteFile(fileURL string) error {
	// Extract key from URL
	fileName := filepath.Base(fileURL)

	// Delete from S3
	_, err := u.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(u.bucketName),
		Key:    aws.String(fileName),
	})

	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}
