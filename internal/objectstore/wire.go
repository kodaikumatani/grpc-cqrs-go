package objectstore

import (
	"context"
	"os"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	ProvideUploader,
)

func ProvideUploader(ctx context.Context) (Uploader, func(), error) {
	backend := os.Getenv("OBJECT_STORE_BACKEND")

	switch backend {
	case "gcs":
		endpoint := os.Getenv("GCS_ENDPOINT")
		u, err := NewGCSUploader(ctx, endpoint)
		if err != nil {
			return nil, nil, err
		}
		return u, func() { u.Close() }, nil
	default:
		endpoint := os.Getenv("MINIO_ENDPOINT")
		accessKey := os.Getenv("MINIO_ACCESS_KEY")
		secretKey := os.Getenv("MINIO_SECRET_KEY")
		useSSL := os.Getenv("MINIO_USE_SSL") == "true"
		u, err := NewMinIOUploader(endpoint, accessKey, secretKey, useSSL)
		if err != nil {
			return nil, nil, err
		}
		return u, func() {}, nil
	}
}
