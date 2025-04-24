package config

import (
	"fmt"
	"os"

	"github.com/latoulicious/siresto-backend/internal/utils"
)

// NewR2UploaderFromEnv creates a new R2Uploader instance using environment variables
func NewR2UploaderFromEnv() (*utils.R2Uploader, error) {
	accessKey := os.Getenv("R2_ACCESS_KEY")
	secretKey := os.Getenv("R2_SECRET_KEY")
	endpoint := os.Getenv("R2_ENDPOINT")
	region := os.Getenv("R2_REGION")
	bucketName := os.Getenv("R2_BUCKET_NAME")
	baseURL := os.Getenv("R2_BASE_URL")

	if accessKey == "" || secretKey == "" || endpoint == "" || region == "" || bucketName == "" || baseURL == "" {
		return nil, fmt.Errorf("missing required Cloudflare R2 configuration")
	}

	return utils.NewR2Uploader(accessKey, secretKey, endpoint, region, bucketName, baseURL)
}
