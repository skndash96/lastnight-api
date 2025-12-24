package repository

import (
	"context"

	"github.com/skndash96/lastnight-backend/internal/db"
)

type uploadRepository struct {
	q *db.Queries
}

func NewUploadRepository(d db.DBTX) *uploadRepository {
	return &uploadRepository{
		q: db.New(d),
	}
}

func (r *uploadRepository) GetOrCreateUpload(ctx context.Context, key, mimeType string, size int64, sha256 string) (*db.GetOrCreateUploadRow, error) {
	upload, err := r.q.GetOrCreateUpload(ctx, db.GetOrCreateUploadParams{
		StorageKey:   key,
		FileMimeType: mimeType,
		FileSize:     size,
		FileSha256:   sha256,
	})
	if err != nil {
		return nil, NewRepoError(err, RepoErrInternal, "Failed to create upload")
	}
	return &upload, nil
}

func (r *uploadRepository) CreateUploadRef(ctx context.Context, uploadID int32, teamID int32, userID int32, name string) (*db.UploadRef, error) {
	ref, err := r.q.CreateUploadRef(ctx, db.CreateUploadRefParams{
		FileName:   name,
		UploadID:   uploadID,
		TeamID:     teamID,
		UploaderID: userID,
	})
	if err != nil {
		return nil, NewRepoError(err, RepoErrInternal, "Failed to create upload reference")
	}
	return &ref, nil
}

func (r *uploadRepository) CreateUploadRefTag(ctx context.Context, uploadRefID int32, keyID int32, valueID int32) error {
	err := r.q.CreateUploadRefTags(ctx, db.CreateUploadRefTagsParams{
		UploadRefID: uploadRefID,
		KeyID:       keyID,
		ValueID:     valueID,
	})

	if err != nil {
		return NewRepoError(err, RepoErrInternal, "Failed to create upload tag")
	}

	return nil
}
