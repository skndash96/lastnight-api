package dto

import "github.com/skndash96/lastnight-backend/internal/db"

// TODO: Refactor DTO so that it does NOT contain any database-specific types

// ------ path params ------
type TeamPathParams struct {
	TeamID int32 `param:"teamID" validate:"required"`
}

type TagPathParams struct {
	TeamPathParams
	TagID int32 `param:"tagID" validate:"required"`
}

type TagValuePathParams struct {
	TagPathParams
	TagValueID int32 `param:"tagValueID" validate:"required"`
}

// ------ body ------
type CreateTagBody struct {
	Name     string         `json:"name" validate:"required,min=2,max=100"`
	DataType db.TagDataType `json:"data_type" validate:"required,min=2,max=100"`
}

type CreateTagValueBody struct {
	Value string `json:"value" validate:"required,min=2,max=100"`
}

type UpdateTagBody struct {
	Name string `json:"name" validate:"required,min=2,max=100"`
	// data type NOT allowed
}

// ------ request ------
type GetTagsRequest struct {
	TeamPathParams
}

type GetTagValuesRequest struct {
	TagPathParams
}

type CreateTagRequest struct {
	TeamPathParams
	CreateTagBody
}

type UpdateTagRequest struct {
	TagPathParams
	UpdateTagBody
}

type DeleteTagRequest struct {
	TagPathParams
}

type CreateTagValueRequest struct {
	TagPathParams
	CreateTagValueBody
}

type DeleteTagValueRequest struct {
	TagValuePathParams
}

// ------ response ------
type GetTagsResponse struct {
	Data []db.Tag `json:"data"`
}

type GetTagValuesResponse struct {
	Data []db.TagValue `json:"data"`
}

type CreateTagResponse struct {
	Data *db.Tag `json:"data"`
}

type UpdateTagResponse struct {
	Data *db.Tag `json:"data"`
}

type DeleteTagResponse struct {
	Data *db.Tag `json:"data"`
}

type CreateTagValueResponse struct {
	Data *db.TagValue `json:"data"`
}

type UpdateTagValueResponse struct {
	Data *db.TagValue `json:"data"`
}

type DeleteTagValueResponse struct {
	Data *db.TagValue `json:"data"`
}
