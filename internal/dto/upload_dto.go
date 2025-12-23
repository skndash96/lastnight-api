package dto

import (
	"github.com/skndash96/lastnight-backend/internal/service"
)

// ------ body ------
type PresignUploadsBody struct {
	Files []*service.PresignUploadItem `json:"files"`
}

type CompleteUploadsBody struct {
	Files []*service.CompleteUploadItem `json:"files"`
}

// ------ request ------
type PresignUploadsRequest struct {
	TeamPathParams
	PresignUploadsBody
}

type CompleteUploadsRequest struct {
	TeamPathParams
	CompleteUploadsBody
}

// ------ response ------
type PresignUploadResult struct {
	Url    string            `json:"url"`
	Fields map[string]string `json:"fields"`
}
type PresignUploadsResponse struct {
	Results []PresignUploadResult `json:"results"`
}
