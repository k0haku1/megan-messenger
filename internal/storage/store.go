package storage

import (
	"context"
	"io"
	"time"
)

// ObjectStore is the abstract object storage used by avatars and attachments.
type ObjectStore interface {
	EnsureBucket(ctx context.Context) error
	Put(ctx context.Context, key string, body io.Reader, size int64, contentType string) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	PresignGet(ctx context.Context, key string, expiry time.Duration) (string, error)
}

const (
	PresignTTL = time.Hour
)
