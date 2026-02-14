package s3

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3 struct {
	client   *minio.Client
	useSSL   bool
	bucket   string
	endpoint string
}

type UploadResponse struct {
	URL      string `json:"url"`
	FileName string `json:"file_name"`
	Size     int64  `json:"size"`
}

func NewS3(endpoint, region, bucket, accessKey, secretKey string, useSSL bool) (*S3, error) {
	ctx := context.Background()

	minioClient, errInit := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if errInit != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", errInit)
	}

	bucketName := bucket

	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: region})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket %s already exists\n", bucketName)
		} else {
			return nil, fmt.Errorf("failed to create/check bucket: %w", err)
		}
	} else {
		log.Printf("Successfully created bucket %s\n", bucketName)
	}

	return &S3{client: minioClient, useSSL: useSSL, bucket: bucket, endpoint: endpoint}, nil
}

func (s *S3) UploadFile(ctx context.Context, file *multipart.FileHeader, folder string) (*UploadResponse, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close() //nolint:errcheck

	fileName := generateUniqueFileName(file.Filename)
	objectName := filepath.Join(folder, fileName)

	info, err := s.client.PutObject(ctx, s.bucket, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	url := s.getFileURL(objectName)

	return &UploadResponse{
		URL:      url,
		FileName: fileName,
		Size:     info.Size,
	}, nil
}

func (s *S3) UploadStream(ctx context.Context, reader io.Reader, objectName, contentType string, size int64) (*UploadResponse, error) {
	info, err := s.client.PutObject(ctx, s.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload stream: %w", err)
	}

	url := s.getFileURL(objectName)

	return &UploadResponse{
		URL:      url,
		FileName: filepath.Base(objectName),
		Size:     info.Size,
	}, nil
}

func (s *S3) DeleteFile(ctx context.Context, objectName string) error {
	err := s.client.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (s *S3) GetFileURL(objectName string) string {
	return s.getFileURL(objectName)
}

func (s *S3) getFileURL(objectName string) string {
	if s.useSSL {
		return fmt.Sprintf("https://%s/%s/%s", s.endpoint, s.bucket, objectName)
	}
	return fmt.Sprintf("http://%s/%s/%s", s.endpoint, s.bucket, objectName)
}

func generateUniqueFileName(originalName string) string {
	ext := filepath.Ext(originalName)
	name := originalName[:len(originalName)-len(ext)]
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s_%d%s", name, timestamp, ext)
}
