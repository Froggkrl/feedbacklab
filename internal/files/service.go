package files

import (
	"context"
	minio_client "innotech/pkg/minio"
	"log/slog"
)

// Service defines the interface for files business logic operations.
type Service struct {
	minio  *minio_client.MinioClient
	logger *slog.Logger
}

// NewService creates a new Service instance.
func NewService(minioClient *minio_client.MinioClient, logger *slog.Logger) *Service {
	return &Service{minio: minioClient, logger: logger}
}

// UploadFile creates a new file to MinIO
func (s *Service) UploadFile(ctx context.Context, filename, path string) (string, error) {
	s.logger.Info("uploading file", "filename", filename)

	err := s.minio.Upload(ctx, filename, path)
	if err != nil {
		s.logger.Error("upload failed", "err", err)
		return "", err
	}

	url, err := s.minio.GetFileURL(filename)
	if err != nil {
		s.logger.Error("cannot get file url", "err", err)
		return "", err
	}

	s.logger.Info("file uploaded", "url", url)
	return url, nil
}
