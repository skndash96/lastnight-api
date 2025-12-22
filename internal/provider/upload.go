package provider

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/skndash96/lastnight-backend/internal/config"
)

type UploadProvider interface {
	PresignUpload(ctx context.Context, teamID int32, file *IncomingFile) (*url.URL, map[string]string, error)
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

type IncomingFile struct {
	Name string
	Size int64
	Mime string
}

func (p *uploadProvider) PresignUpload(ctx context.Context, teamID int32, file *IncomingFile) (*url.URL, map[string]string, error) {
	if file.Size > 11*p.cfg.MaxSize/10 {
		return nil, nil, ErrFileTooLarge
	} else if !strings.HasPrefix(file.Mime, "image/") && !strings.HasPrefix(file.Mime, "application/") {
		return nil, nil, ErrInvalidFileType
	}

	policy := minio.NewPostPolicy()
	policy.SetExpires(time.Now().Add(p.cfg.Expiration))
	policy.SetBucket(p.cfg.BucketName)
	policy.SetContentLengthRange(9*file.Size/10, 11*file.Size/10)
	policy.SetContentType(file.Mime)
	policy.SetKey(p.GenerateObjectKey(strconv.FormatInt(int64(teamID), 10), file.Name))

	url, fields, err := p.client.PresignedPostPolicy(ctx, policy)
	if err != nil {
		return nil, nil, err
	}

	return url, fields, nil
}

func (p *uploadProvider) GenerateObjectKey(teamID, originalName string) string {
	ext := path.Ext(originalName)
	if len(ext) > 10 { // sanity cap
		ext = ""
	}

	id := uuid.NewString()

	return fmt.Sprintf("team_%s/uploads/%s%s", teamID, id, strings.ToLower(ext))
}
