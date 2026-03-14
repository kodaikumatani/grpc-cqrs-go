package objectstore

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOUploader struct {
	client *minio.Client
}

func NewMinIOUploader(endpoint, accessKey, secretKey string, useSSL bool) (*MinIOUploader, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	return &MinIOUploader{client: client}, nil
}

func (u *MinIOUploader) Upload(ctx context.Context, bucket, key string, r io.Reader) error {
	_, err := u.client.PutObject(ctx, bucket, key, r, -1, minio.PutObjectOptions{})
	return err
}
