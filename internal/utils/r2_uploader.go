package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type R2Uploader struct {
	Client     *s3.Client
	BucketName string
	BaseURL    string
}

func NewR2Uploader(accessKey, secretKey, endpoint, region, bucketName, baseURL string) (*R2Uploader, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
		config.WithRegion(region),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: endpoint, HostnameImmutable: true}, nil
			}),
		),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	uploader := &R2Uploader{
		Client:     client,
		BucketName: bucketName,
		BaseURL:    baseURL,
	}

	// Configure CORS for the bucket
	if err := uploader.ConfigureCORS(); err != nil {
		log.Printf("Warning: Failed to configure CORS for bucket %s: %v", bucketName, err)
		// Continue anyway - don't fail initialization
	}

	return uploader, nil
}

// ConfigureCORS sets up CORS configuration for the R2 bucket
// This configuration applies to all objects in the bucket, both existing and new uploads
func (u *R2Uploader) ConfigureCORS() error {
	corsConfig := &s3.PutBucketCorsInput{
		Bucket: aws.String(u.BucketName),
		CORSConfiguration: &types.CORSConfiguration{
			CORSRules: []types.CORSRule{
				{
					AllowedHeaders: []string{"*"},
					AllowedMethods: []string{"GET", "HEAD", "PUT", "POST", "DELETE"},
					// Allow the specific frontend domain instead of wildcard
					AllowedOrigins: []string{
						"http://localhost:41234",
						"http://127.0.0.1:41234",
						"*", // Keep wildcard for testing but remove in production
					},
					ExposeHeaders: []string{"ETag"},
					MaxAgeSeconds: aws.Int32(3600),
				},
			},
		},
	}

	_, err := u.Client.PutBucketCors(context.TODO(), corsConfig)
	if err != nil {
		return fmt.Errorf("failed to set CORS configuration: %w", err)
	}

	log.Printf("Successfully configured CORS for bucket: %s", u.BucketName)
	return nil
}

func (u *R2Uploader) Upload(file io.Reader, filename string) (string, error) {
	buffer := new(bytes.Buffer)
	if _, err := io.Copy(buffer, file); err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}

	// Detect content type
	fileBytes := buffer.Bytes()
	contentType := http.DetectContentType(fileBytes)

	// Use extension as fallback for certain image types
	ext := filepath.Ext(filename)
	if contentType == "application/octet-stream" {
		switch ext {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".gif":
			contentType = "image/gif"
		case ".webp":
			contentType = "image/webp"
		case ".svg":
			contentType = "image/svg+xml"
		}
	}

	uploader := manager.NewUploader(u.Client)
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(u.BucketName),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(fileBytes),
		ACL:         "public-read",
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("upload to R2: %w", err)
	}

	return fmt.Sprintf("%s/%s", u.BaseURL, filename), nil
}
