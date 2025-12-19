package repository

import (
	"context"
	"encoding/json"

	"github.com/skndash96/lastnight-backend/internal/db"
)

type TagRepo struct {
	q *db.Queries
}

func NewTagRepo(d db.DBTX) *TagRepo {
	return &TagRepo{
		q: db.New(d),
	}
}

func (r *TagRepo) ListFilters(ctx context.Context, membershipID int32) ([]db.Tag, error) {
	raw, err := r.q.ListFilters(ctx, membershipID)
	if err != nil {
		return nil, err
	}

	out := []db.Tag{}
	for _, t := range raw {
		var tag db.Tag
		tag.KeyID = t.KeyID
		tag.Key = t.Key
		if t.ValueID.Valid {
			tag.ValueID = t.ValueID.Int32
		}
		if t.Value.Valid {
			tag.Value = t.Value.String
		}

		err := json.Unmarshal(t.Options, &tag.Options)
		if err != nil {
			return nil, err
		}
		out = append(out, tag)
	}

	return out, nil
}

func (r *TagRepo) CreateTagKey(ctx context.Context, teamID int32, tagName string, dataType db.TagDataType) (db.TagKey, error) {
	return r.q.CreateTagKey(ctx, db.CreateTagKeyParams{
		TeamID:   teamID,
		Name:     tagName,
		DataType: dataType,
	})
}

func (r *TagRepo) CreateTagValue(ctx context.Context, tagID int32, value string) (db.TagValue, error) {
	return r.q.CreateTagValue(ctx, db.CreateTagValueParams{
		KeyID: tagID,
		Value: value,
	})
}

func (r *TagRepo) DeleteTagKey(ctx context.Context, tagID int32) (db.TagKey, error) {
	return r.q.DeleteTagKey(ctx, tagID)
}

func (r *TagRepo) DeleteTagValue(ctx context.Context, tagValueID int32) (db.TagValue, error) {
	return r.q.DeleteTagValue(ctx, tagValueID)
}

func (r *TagRepo) UpdateTagKey(ctx context.Context, tagID int32, tagName string) (db.TagKey, error) {
	return r.q.UpdateTagKey(ctx, db.UpdateTagKeyParams{
		ID:   tagID,
		Name: tagName,
	})
}
