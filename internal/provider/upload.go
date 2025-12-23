package provider

import (
	"context"
	"errors"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/skndash96/lastnight-backend/internal/config"
)

type UploadProvider interface {
	PresignUpload(ctx context.Context, key, name, mime string, size int64) (*url.URL, map[string]string, error)
	MoveObject(ctx context.Context, dstKey, srcKey string) error
}

type uploadProvider struct {
	client *minio.Client
	cfg    *config.MinioConfig
}

var (
	ErrFileTooLarge    = errors.New("file too large")
	ErrInvalidFileType = errors.New("invalid file type")
)

func NewUploadProvider(cfg config.MinioConfig) (UploadProvider, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Username, cfg.Password, ""),
		Secure: cfg.UseSSL,
	})

	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	err = client.MakeBucket(ctx, cfg.BucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := client.BucketExists(ctx, cfg.BucketName)
		if errBucketExists != nil || !exists {
			return nil, err
		}
	}

	return &uploadProvider{
		client: client,
		cfg:    &cfg,
	}, nil
}

func (p *uploadProvider) PresignUpload(ctx context.Context, key, name, mime string, size int64) (*url.URL, map[string]string, error) {
	if size > p.cfg.MaxSize {
		return nil, nil, ErrFileTooLarge
	}

	policy := minio.NewPostPolicy()
	policy.SetExpires(time.Now().Add(p.cfg.Expiration))
	policy.SetBucket(p.cfg.BucketName)
	policy.SetContentLengthRange(9*size/10, 11*size/10)
	policy.SetContentType(mime)
	policy.SetKey(key)

	url, fields, err := p.client.PresignedPostPolicy(ctx, policy)
	if err != nil {
		return nil, nil, err
	}

	return url, fields, nil
}

func (p *uploadProvider) MoveObject(ctx context.Context, dst, src string) error {
	_, err := p.client.CopyObject(ctx, minio.CopyDestOptions{
		Bucket: p.cfg.BucketName,
		Object: dst,
	}, minio.CopySrcOptions{
		Bucket: p.cfg.BucketName,
		Object: src,
	})

	if err != nil {
		return err
	}

	err = p.client.RemoveObject(ctx, p.cfg.BucketName, src, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (p *uploadProvider) DeleteObject(ctx context.Context, objectKey string) error {
	err := p.client.RemoveObject(ctx, p.cfg.BucketName, objectKey, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
