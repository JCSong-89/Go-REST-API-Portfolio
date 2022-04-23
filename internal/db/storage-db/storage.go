package storage_db

import (
	"github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/credentials"
)

func NewStorageClient(endpoint string, accessKeyID string, secretAccessKey string, useSSL bool) (*minio.Client, error) {
	minioClient, err := minio.NewWithOptions(endpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		})

	if err != nil {
		return nil, err
	}

	return minioClient, nil
}