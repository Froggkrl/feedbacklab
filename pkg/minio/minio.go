// Package minio - package to create basic operation in MinIO storage
package minio

import (
	"context"
	"log"
	"net/url"
	"time"

	"innotech/pkg/logger"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient provides helper methods for interacting with MinIO storage.
type MinioClient struct {
	Client     *minio.Client
	BucketName string
}

func New(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinioClient, error) {
	logger.Info("initializing MinIO client",
		"endpoint", endpoint,
		"bucket", bucket,
		"useSSL", useSSL,
	)

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logger.Error("failed to create MinIO client",
			"endpoint", endpoint,
			"error", err,
		)
		return nil, err
	}

	ctx := context.Background()
	exists, _ := client.BucketExists(ctx, bucket)
	if !exists {
		logger.Info("bucket does not exist, creating new bucket",
			"bucket", bucket,
		)

		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			logger.Error("failed to create bucket",
				"bucket", bucket,
				"error", err,
			)
			return nil, err
		}
		log.Println("Bucket created:", bucket)
		logger.Info("bucket created successfully",
			"bucket", bucket,
		)
	} else {
		logger.Debug("bucket already exists",
			"bucket", bucket,
		)
	}

	logger.Info("MinIO client initialized successfully",
		"endpoint", endpoint,
		"bucket", bucket,
	)
	return &MinioClient{
		Client:     client,
		BucketName: bucket,
	}, nil
}

func (m *MinioClient) Upload(ctx context.Context, objectName, filePath string) error {
	logger.Info("uploading file to MinIO",
		"bucket", m.BucketName,
		"object_name", objectName,
		"file_path", filePath,
	)

	_, err := m.Client.FPutObject(ctx, m.BucketName, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		logger.Error("failed to upload file to MinIO",
			"bucket", m.BucketName,
			"object_name", objectName,
			"file_path", filePath,
			"error", err,
		)
		return err
	}

	logger.Info("file uploaded successfully",
		"bucket", m.BucketName,
		"object_name", objectName,
		"file_path", filePath,
	)
	return nil
}

func (m *MinioClient) GetFileURL(objectName string) (string, error) {
	logger.Debug("generating presigned URL for object",
		"bucket", m.BucketName,
		"object_name", objectName,
	)

	reqParams := make(url.Values)

	presignedURL, err := m.Client.PresignedGetObject(
		context.Background(),
		m.BucketName,
		objectName,
		time.Hour*24,
		reqParams,
	)
	if err != nil {
		logger.Error("failed to generate presigned URL",
			"bucket", m.BucketName,
			"object_name", objectName,
			"error", err,
		)
		return "", err
	}

	urlString := presignedURL.String()
	logger.Debug("presigned URL generated successfully",
		"bucket", m.BucketName,
		"object_name", objectName,
		"url_length", len(urlString),
	)

	return urlString, nil
}
