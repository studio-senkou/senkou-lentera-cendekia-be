package storage

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/studio-senkou/lentera-cendekia-be/utils/app"
)

type S3Config struct {
	Host            string `json:"host"`
	Region          string `json:"region"`
	Bucket          string `json:"bucket"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}

func NewS3Config(host, region, bucket, accessKeyID, secretAccessKey string) *S3Config {
	return &S3Config{
		Host:            host,
		Region:          region,
		Bucket:          bucket,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}
}

func NewStorage(config *S3Config) *aws.Config {
	return &aws.Config{
		Credentials: credentials.NewStaticCredentialsProvider(
			config.AccessKeyID,
			config.SecretAccessKey,
			"",
		),
		Region:       config.Region,
		BaseEndpoint: aws.String(config.Host),
	}
}

func GetStorage() *aws.Config {
	config := NewS3Config(
		app.GetEnv("AWS_S3_HOST", "https://is3.cloudhost.id"),
		app.GetEnv("AWS_S3_REGION", "sgp01"),
		app.GetEnv("AWS_S3_BUCKET", "lentera-cendekia"),
		app.GetEnv("AWS_S3_ACCESS_KEY_ID", ""),
		app.GetEnv("AWS_S3_SECRET_ACCESS_KEY", ""),
	)

	return NewStorage(config)
}

func GetBucket() string {
	return app.GetEnv("AWS_S3_BUCKET", "lentera-cendekia")
}

type UploadService struct {
	storage *s3.Client
	bucket  string
}

func NewUploadService() *UploadService {
	config := GetStorage()
	bucket := GetBucket()

	return &UploadService{
		storage: s3.NewFromConfig(*config),
		bucket:  bucket,
	}
}

func (s *UploadService) UploadFile(ctx context.Context, file *multipart.FileHeader, filename string, folder *string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	var key string
	if folder != nil && *folder != "" {
		key = fmt.Sprintf("%s/%s", *folder, filename)
	} else {
		key = filename
	}

	_, err = s.storage.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          src,
		ContentType:   aws.String(file.Header.Get("Content-Type")),
		ContentLength: aws.Int64(file.Size),
		ACL:           types.ObjectCannedACLPublicRead,
		CacheControl:  aws.String("max-age=31536000, public"),
		Metadata: map[string]string{
			"filename":    file.Filename,
			"uploaded-by": "Lentera Cendekia API",
		},
	})
	if err != nil {
		return "", err
	}

	return key, nil
}

func (s *UploadService) RemoveFile(ctx context.Context, path string) error {
	_, err := s.storage.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return fmt.Errorf("failed to remove file: %w", err)
	}
	return nil
}
