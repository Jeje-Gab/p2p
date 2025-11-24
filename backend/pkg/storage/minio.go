package storage

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/cs2-p2p-skins/backend/pkg/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOStorage struct {
	client    *minio.Client
	bucket    string
	publicURL string
}

func NewMinIOStorage(cfg config.MinIOConfig) (*MinIOStorage, error) {
	// Initialize MinIO client
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	storage := &MinIOStorage{
		client:    client,
		bucket:    cfg.Bucket,
		publicURL: cfg.PublicURL,
	}

	// Create bucket if it doesn't exist
	if err := storage.ensureBucket(context.Background()); err != nil {
		return nil, err
	}

	log.Printf("[MINIO] Connected to %s | Bucket: %s", cfg.Endpoint, cfg.Bucket)
	return storage, nil
}

// ensureBucket creates the bucket if it doesn't exist and sets it to public
func (s *MinIOStorage) ensureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		log.Printf("[MINIO] Creating bucket: %s", s.bucket)
		if err := s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("[MINIO] Bucket created successfully: %s", s.bucket)
	}

	// Set bucket policy to public read
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, s.bucket)

	if err := s.client.SetBucketPolicy(ctx, s.bucket, policy); err != nil {
		log.Printf("[MINIO] Warning: failed to set bucket policy to public: %v", err)
	} else {
		log.Printf("[MINIO] Bucket policy set to public: %s", s.bucket)
	}

	return nil
}

// UploadFile uploads a file to MinIO
func (s *MinIOStorage) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Return public URL
	url := fmt.Sprintf("%s/%s/%s", s.publicURL, s.bucket, objectName)
	return url, nil
}

// DownloadFile downloads a file from a URL and uploads it to MinIO
func (s *MinIOStorage) DownloadAndUpload(ctx context.Context, sourceURL, objectName string) (string, error) {
	// This will be implemented in the migration script
	return "", fmt.Errorf("not implemented")
}

// GetPublicURL returns the public URL for an object
func (s *MinIOStorage) GetPublicURL(objectName string) string {
	return fmt.Sprintf("%s/%s/%s", s.publicURL, s.bucket, objectName)
}

// FileExists checks if a file exists in the bucket
func (s *MinIOStorage) FileExists(ctx context.Context, objectName string) (bool, error) {
	_, err := s.client.StatObject(ctx, s.bucket, objectName, minio.StatObjectOptions{})
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
