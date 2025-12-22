package dto

import (
	"github.com/skndash96/lastnight-backend/internal/provider"
)

// ------ body ------
type PresignUploadsBody struct {
	Files []*provider.IncomingFile `json:"files"`
}

// ------ request ------
type PresignUploadsRequest struct {
	TeamPathParams
	PresignUploadsBody
}

// ------ response ------
type PresignUploadResult struct {
	Url    string            `json:"url"`
	Fields map[string]string `json:"fields"`
}
type PresignUploadsResponse struct {
	Results []PresignUploadResult `json:"results"`
}
