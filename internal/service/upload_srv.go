package service

import (
	"context"
	"fmt"
	"net/url"

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

func (s *UploadService) PresignUploads(ctx context.Context, teamID int32, files []*provider.IncomingFile) ([]PresignUploadResult, error) {
	results := make([]PresignUploadResult, len(files))

	for i, file := range files {
		url, fields, err := s.uploadProvider.PresignUpload(ctx, teamID, file)
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
