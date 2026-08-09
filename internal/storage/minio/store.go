package miniostore

import (
	"context"
	"fmt"
	"io"
	"megan-messenger/internal/config"
	"megan-messenger/internal/storage"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Store struct {
	client     *minio.Client
	bucket     string
	publicBase string
}

func New(cfg config.MinIOConfig) (*Store, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio client: %w", err)
	}

	return &Store{
		client:     client,
		bucket:     cfg.Bucket,
		publicBase: strings.TrimRight(cfg.PublicBase, "/"),
	}, nil
}

func (s *Store) EnsureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return fmt.Errorf("bucket exists: %w", err)
	}
	if !exists {
		if err := s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("make bucket: %w", err)
		}
	}

	// Public read for MVP so browser can load via PublicBase (localhost:9000).
	policy := fmt.Sprintf(`{
	  "Version":"2012-10-17",
	  "Statement":[{
	    "Effect":"Allow",
	    "Principal":{"AWS":["*"]},
	    "Action":["s3:GetObject"],
	    "Resource":["arn:aws:s3:::%s/*"]
	  }]
	}`, s.bucket)
	if err := s.client.SetBucketPolicy(ctx, s.bucket, policy); err != nil {
		return fmt.Errorf("set bucket policy: %w", err)
	}
	return nil
}

func (s *Store) Put(ctx context.Context, key string, body io.Reader, size int64, contentType string) error {
	opts := minio.PutObjectOptions{}
	if contentType != "" {
		opts.ContentType = contentType
	}
	_, err := s.client.PutObject(ctx, s.bucket, key, body, size, opts)
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}
	return nil
}

func (s *Store) Delete(ctx context.Context, key string) error {
	if key == "" {
		return nil
	}
	err := s.client.RemoveObject(ctx, s.bucket, key, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("remove object: %w", err)
	}
	return nil
}

func (s *Store) Exists(ctx context.Context, key string) (bool, error) {
	if key == "" {
		return false, nil
	}
	_, err := s.client.StatObject(ctx, s.bucket, key, minio.StatObjectOptions{})
	if err != nil {
		errResp := minio.ToErrorResponse(err)
		if errResp.Code == "NoSuchKey" || errResp.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Store) PresignGet(ctx context.Context, key string, expiry time.Duration) (string, error) {
	if key == "" {
		return "", nil
	}
	if s.publicBase != "" {
		return s.publicBase + "/" + strings.TrimLeft(key, "/"), nil
	}
	if expiry <= 0 {
		expiry = storage.PresignTTL
	}
	u, err := s.client.PresignedGetObject(ctx, s.bucket, key, expiry, url.Values{})
	if err != nil {
		return "", fmt.Errorf("presign get: %w", err)
	}
	return u.String(), nil
}

var _ storage.ObjectStore = (*Store)(nil)
