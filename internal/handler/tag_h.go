package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/skndash96/lastnight-backend/internal/auth"
	"github.com/skndash96/lastnight-backend/internal/dto"
	"github.com/skndash96/lastnight-backend/internal/service"
)

type tagHandler struct {
	tagSrv *service.TagService
}

func NewTagHandler(s *service.TagService) *tagHandler {
	return &tagHandler{
		tagSrv: s,
	}
}

// GetTags retrieves the tag of a team.
// @Summary Get Tag
// @Tags Tag
// @Description Get the tags of a team
// @Param teamID path string true "Team ID"
// @Produce json
// @Success 200 {object} dto.GetTagsResponse
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/{teamID}/tags [get]
func (h *tagHandler) ListTags(c echo.Context) error {
	v := new(dto.GetTagsRequest)
	if err := c.Bind(v); err != nil {
		return err
	}

	session, ok := auth.GetSession(c)
	if !ok {
		return echo.ErrUnauthorized
	}

	tags, err := h.tagSrv.ListTags(c.Request().Context(), session.MembershipID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.GetTagsResponse{
		Data: tags,
	})
}

// @Summary New Tag Key
// @Tags Tag
// @Description Create a new tag key
// @Param teamID path string true "Team ID"
// @Param tag body dto.CreateTagKeyBody true "Tag"
// @Produce json
// @Success 201 {object} dto.CreateTagKeyResponse
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/{teamID}/tags [post]
func (h *tagHandler) CreateTagKey(c echo.Context) error {
	v := new(dto.CreateTagKeyRequest)
	if err := c.Bind(v); err != nil {
		return err
	}

	if err := c.Validate(v); err != nil {
		return err
	}

	tag, err := h.tagSrv.CreateTagKey(c.Request().Context(), v.TeamID, v.Name, v.DataType)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, dto.CreateTagKeyResponse{
		Data: tag,
	})
}

// @Summary Update Tag Key
// @Tags Tag
// @Description Update a tag key
// @Param teamID path string true "Team ID"
// @Param tagID path string true "Tag ID"
// @Param tag body dto.UpdateTagKeyBody true "Tag"
// @Produce json
// @Success 200 {object} dto.UpdateTagKeyResponse
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/{teamID}/tags/{tagID} [put]
func (h *tagHandler) UpdateTagKey(c echo.Context) error {
	v := new(dto.UpdateTagKeyRequest)
	if err := c.Bind(v); err != nil {
		return err
	}

	if err := c.Validate(v); err != nil {
		return err
	}

	tag, err := h.tagSrv.UpdateTagKey(c.Request().Context(), v.TagID, v.Name)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.UpdateTagKeyResponse{
		Data: tag,
	})
}

// @Summary Delete Tag Key
// @Tags Tag
// @Description Delete a tag key
// @Param teamID path string true "Team ID"
// @Param tagID path string true "Tag ID"
// @Produce json
// @Success 200 {object} dto.DeleteTagKeyResponse
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/{teamID}/tags/{tagID} [delete]
func (h *tagHandler) DeleteTagKey(c echo.Context) error {
	v := new(dto.DeleteTagKeyRequest)

	if err := c.Bind(v); err != nil {
		return err
	}

	if err := c.Validate(v); err != nil {
		return err
	}

	tag, err := h.tagSrv.DeleteTagKey(c.Request().Context(), v.TagID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.DeleteTagKeyResponse{
		Data: tag,
	})
}

// @Summary Create Tag Value
// @Tags Tag
// @Description Create a new tag value
// @Param teamID path string true "Team ID"
// @Param tagID path string true "Tag ID"
// @Param value body dto.CreateTagValueBody true "Value"
// @Produce json
// @Success 201 {object} dto.CreateTagValueResponse
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/{teamID}/tags/{tagID}/values [post]
func (h *tagHandler) CreateTagValue(c echo.Context) error {
	v := new(dto.CreateTagValueRequest)

	if err := c.Bind(v); err != nil {
		return err
	}

	if err := c.Validate(v); err != nil {
		return err
	}

	value, err := h.tagSrv.CreateTagValue(c.Request().Context(), v.TagID, v.Value)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, dto.CreateTagValueResponse{
		Data: value,
	})
}

// @Summary Delete Tag Value
// @Tags Tag
// @Description Delete a tag value
// @Param teamID path string true "Team ID"
// @Param tagID path string true "Tag ID"
// @Param valueID path string true "Value ID"
// @Produce json
// @Success 200 {object} dto.DeleteTagValueResponse
// @Failure default {object} dto.ErrorResponse
// @Router /api/teams/{teamID}/tags/{tagID}/values/{valueID} [delete]
func (h *tagHandler) DeleteTagValue(c echo.Context) error {
	v := new(dto.DeleteTagValueRequest)

	if err := c.Bind(v); err != nil {
		return err
	}

	if err := c.Validate(v); err != nil {
		return err
	}

	value, err := h.tagSrv.DeleteTagValue(c.Request().Context(), v.TagValueID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.DeleteTagValueResponse{
		Data: value,
	})
}
