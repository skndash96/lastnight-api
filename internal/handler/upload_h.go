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
// @Param upload_request body dto.PresignUploadBody true "Presign request"
// @Produce json
// @Success 201 {object} dto.PresignUploadResponse
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/{teamID}/uploads/presign [post]
func (h *uploadHandler) PresignUpload(c echo.Context) error {
	session, ok := auth.GetSession(c)
	if !ok {
		return echo.ErrUnauthorized
	}

	v := dto.PresignUploadRequest{}
	if err := c.Bind(&v); err != nil {
		return err
	}

	if err := c.Validate(&v); err != nil {
		return err
	}

	result, err := h.uploadSrv.PresignUpload(c.Request().Context(), session.TeamID, v.Name, v.MimeType, v.Size)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, &dto.PresignUploadResponse{
		Url:    result.Url.String(),
		Fields: result.Fields,
	})

	return nil
}

// @Summary Complete upload
// @Description Call this route after client-side uploading to the bucket via POST policy. Processes uploaded file.
// @Tags Upload
// @Accept json
// @Param teamID path string true "Team ID"
// @Param upload_request body dto.CompleteUploadBody true "Complete upload request"
// @Produce json
// @Success 200
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/{teamID}/uploads/complete [post]
func (h *uploadHandler) CompleteUpload(c echo.Context) error {
	session, ok := auth.GetSession(c)
	if !ok {
		return echo.ErrUnauthorized
	}

	v := dto.CompleteUploadRequest{}
	if err := c.Bind(&v); err != nil {
		return err
	}

	if err := c.Validate(&v); err != nil {
		return err
	}

	tags := make([][]int32, len(v.Tags))
	for i, tag := range v.Tags {
		tags[i] = []int32{tag.KeyID, tag.ValueID}
	}

	err := h.uploadSrv.CompleteUpload(c.Request().Context(), session.TeamID, session.UserID, v.Key, v.Name, v.MimeType, tags)
	if err != nil {
		return err
	}

	c.NoContent(http.StatusCreated)

	return nil
}
