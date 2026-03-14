package objectstore

import (
	"context"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCSUploader struct {
	client *storage.Client
}

func NewGCSUploader(ctx context.Context, endpoint string) (*GCSUploader, error) {
	opts := []option.ClientOption{
		option.WithoutAuthentication(),
	}
	if endpoint != "" {
		opts = append(opts, option.WithEndpoint(endpoint))
	}

	client, err := storage.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &GCSUploader{client: client}, nil
}

func (u *GCSUploader) Upload(ctx context.Context, bucket, key string, r io.Reader) error {
	w := u.client.Bucket(bucket).Object(key).NewWriter(ctx)
	if _, err := io.Copy(w, r); err != nil {
		w.Close()
		return err
	}
	return w.Close()
}

func (u *GCSUploader) Close() error {
	return u.client.Close()
}
