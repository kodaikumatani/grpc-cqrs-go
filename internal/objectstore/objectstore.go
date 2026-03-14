package objectstore

import (
	"context"
	"io"
)

// Uploader はオブジェクトストレージへのアップロードを抽象化するインターフェース。
type Uploader interface {
	Upload(ctx context.Context, bucket, key string, r io.Reader) error
}
