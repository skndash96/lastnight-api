package service

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/skndash96/lastnight-backend/internal/provider"
)

type UploadService struct {
	uploadProvider provider.UploadProvider
}

func NewUploadService(uploadProvider provider.UploadProvider) *UploadService {
	return &UploadService{
		uploadProvider: uploadProvider,
	}
}

type PresignUploadResult struct {
	Url    *url.URL
	Fields map[string]string
}

type PresignUploadItem struct {
	Name string
	Mime string
	Size int64
}

func (s *UploadService) PresignUploads(ctx context.Context, teamID int32, files []*PresignUploadItem) ([]PresignUploadResult, error) {
	results := make([]PresignUploadResult, len(files))

	for _, file := range files {
		// Allow only application/* mime types
		// TODO: Support image/* mime types, also handle combining images to PDF if done so
		if !strings.HasPrefix(file.Mime, "application/") {
			return nil, NewSrvError(nil, SrvErrInvalidInput, "Upload presign failed: Only application/* mime types are allowed.")
		}
	}

	for i, file := range files {
		key := generateTmpObjectKey(teamID, file.Name)
		url, fields, err := s.uploadProvider.PresignUpload(ctx, key, file.Name, file.Mime, file.Size)
		if err != nil {
			// TODO: Handle error properly
			return nil, NewSrvError(err, SrvErrInternal, fmt.Sprintf("failed to presign upload for file %s", file.Name))
		}

		results[i] = PresignUploadResult{
			Url:    url,
			Fields: fields,
		}
	}

	return results, nil
}

type CompleteUploadItem struct {
	Key  string
	Name string
	Mime string
	Size int64
}

func (s *UploadService) CompleteUploads(ctx context.Context, files []*CompleteUploadItem) error {
	for _, file := range files {
		src := file.Key
		// TODO: Validate file properly, for now just skip
		dst := convertTmpKey(src)
		if dst == "" {
			return NewSrvError(nil, SrvErrInvalidInput, "Upload completion failed: invalid key")
		}

		err := s.uploadProvider.MoveObject(ctx, dst, src)
		if err != nil {
			return NewSrvError(err, SrvErrInternal, fmt.Sprintf("failed to move object %s to %s", src, dst))
		}

		// Add to upload repo
		fmt.Printf("Adding file %s to upload repository\n%v\n", dst, file)
		// Push to queue
	}

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
