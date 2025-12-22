package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/skndash96/lastnight-backend/internal/auth"
	"github.com/skndash96/lastnight-backend/internal/dto"
	"github.com/skndash96/lastnight-backend/internal/service"
)

type uploadHandler struct {
	uploadSrv *service.UploadService
}

func NewUploadHandler(uploadSrv *service.UploadService) *uploadHandler {
	return &uploadHandler{
		uploadSrv: uploadSrv,
	}
}

// @Summary Create pre-signed request
// @Description Create a pre-signed request for uploading files to S3 via POST policy
// @Tags Upload
// @Accept json
// @Param teamID path string true "Team ID"
// @Param upload_request body dto.PresignUploadsBody true "Presign request"
// @Produce json
// @Success 201 {object} dto.PresignUploadsResponse
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/{teamID}/uploads/presign [post]
func (h *uploadHandler) PresignUploads(c echo.Context) error {
	session, ok := auth.GetSession(c)
	if !ok {
		return echo.ErrUnauthorized
	}

	v := dto.PresignUploadsRequest{}
	if err := c.Bind(&v); err != nil {
		return err
	}

	if err := c.Validate(&v); err != nil {
		return err
	}

	_results, err := h.uploadSrv.PresignUploads(c.Request().Context(), session.TeamID, v.Files)
	if err != nil {
		return err
	}

	results := make([]dto.PresignUploadResult, len(_results))
	for i, result := range _results {
		results[i] = dto.PresignUploadResult{
			Url:    result.Url.String(),
			Fields: result.Fields,
		}
	}

	c.JSON(http.StatusOK, &dto.PresignUploadsResponse{
		Results: results,
	})

	return nil
}
