package service

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skndash96/lastnight-backend/internal/provider"
	"github.com/skndash96/lastnight-backend/internal/repository"
)

type UploadService struct {
	pool           *pgxpool.Pool
	uploadProvider provider.UploadProvider
}

func NewUploadService(uploadProvider provider.UploadProvider, pool *pgxpool.Pool) *UploadService {
	return &UploadService{
		pool:           pool,
		uploadProvider: uploadProvider,
	}
}

type PresignUploadResult struct {
	Url    *url.URL
	Fields map[string]string
}

func (s *UploadService) PresignUpload(ctx context.Context, teamID int32, name, mimeType string, size int64) (*PresignUploadResult, error) {
	// Allow only application/* mime types
	// TODO: Support image/* mime types while combining images to PDF if done so
	if !strings.HasPrefix(mimeType, "application/") {
		return nil, NewSrvError(nil, SrvErrInvalidInput, "Upload presign failed: Only application/* mime types are allowed.")
	}

	key := generateTmpObjectKey(teamID, name)
	url, fields, err := s.uploadProvider.PresignUpload(ctx, key, name, mimeType, size)
	if err != nil {
		// TODO: Handle error properly
		return nil, NewSrvError(err, SrvErrInternal, fmt.Sprintf("failed to presign upload for file %s", name))
	}

	return &PresignUploadResult{
		Url:    url,
		Fields: fields,
	}, nil
}

func (s *UploadService) CompleteUpload(ctx context.Context, teamID, userID int32, tmpKey, name, mime string, tags [][]int32) error {
	info, err := s.uploadProvider.GetUploadInfo(ctx, tmpKey)
	if err != nil {
		return NewSrvError(err, SrvErrInternal, fmt.Sprintf("failed to get upload info for %s", tmpKey))
	}

	newKey := convertTmpKey(tmpKey)
	if newKey == "" {
		return NewSrvError(nil, SrvErrInvalidInput, "Upload completion failed: invalid key")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return NewSrvError(err, SrvErrInternal, "failed to begin transaction")
	}
	defer tx.Rollback(ctx)

	uploadRepo := repository.NewUploadRepository(tx)

	// (hash, size) duplication check happens here
	upload, err := uploadRepo.GetOrCreateUpload(ctx, newKey, mime, info.Size, info.SHA256)
	if err != nil {
		return NewSrvError(err, SrvErrInternal, fmt.Sprintf("failed to create upload for %s", newKey))
	}

	uploadRef, err := uploadRepo.CreateUploadRef(ctx, upload.ID, teamID, userID, name)
	if err != nil {
		return NewSrvError(err, SrvErrInternal, fmt.Sprintf("failed to create upload reference for %s", newKey))
	}

	for _, tag := range tags {
		if err := uploadRepo.CreateUploadRefTag(ctx, uploadRef.ID, tag[0], tag[1]); err != nil {
			return NewSrvError(err, SrvErrInternal, fmt.Sprintf("failed to create upload tag for %s", newKey))
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return NewSrvError(err, SrvErrInternal, "failed to complete upload")
	}

	if upload.Created == false {
		fmt.Printf("Upload already exists, using duplicate: %s\n", upload.StorageKey)
		if err := s.uploadProvider.DeleteObject(ctx, tmpKey); err != nil {
			// it's not a fatal error, so we can continue
			fmt.Printf("failed to delete duplicate upload %s: %v\n", tmpKey, err)
		}
	} else {
		// TODO: Retry worker
		err = s.uploadProvider.MoveObject(ctx, newKey, tmpKey)
		if err != nil {
			fmt.Printf("FATAL: failed to move upload from %s to %s: %v\n", tmpKey, newKey, err)
		}
	}

	// TODO: push to queue

	return nil
}

func generateTmpObjectKey(teamID int32, originalName string) string {
	ext := path.Ext(originalName)
	if len(ext) > 10 {
		ext = ""
	}

	id := uuid.NewString()

	return fmt.Sprintf("tmp/team_%d/uploads/%s%s", teamID, id, strings.ToLower(ext))
}

func convertTmpKey(key string) string {
	if strings.HasPrefix(key, "tmp/") {
		return strings.Replace(key, "tmp/", "", 1)
	}
	return ""
}
