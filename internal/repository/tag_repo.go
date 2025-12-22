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
		return nil, NewRepoError(err, RepoErrInternal, "failed to list filters")
	}

	out := []db.Tag{}
	for _, t := range raw {
		var tag db.Tag
		tag.KeyID = t.KeyID
		tag.Key = t.Key
		if t.ValueID.Valid {
			tag.ValueID = &t.ValueID.Int32
		}
		if t.Value.Valid {
			tag.Value = &t.Value.String
		}

		err := json.Unmarshal(t.Options, &tag.Options)
		if err != nil {
			return nil, NewRepoError(err, RepoErrInternal, "failed to unmarshal options")
		}
		out = append(out, tag)
	}

	return out, nil
}

func (r *TagRepo) CreateFilter(ctx context.Context, membershipID int32, tagID int32, tagValueID int32) error {
	err := r.q.CreateFilter(ctx, db.CreateFilterParams{
		MembershipID: membershipID,
		KeyID:        tagID,
		ValueID:      tagValueID,
	})
	if err != nil {
		return NewRepoError(err, RepoErrInternal, "failed to create filter")
	}
	return nil
}

func (r *TagRepo) DeleteAllFilters(ctx context.Context, membershipID int32) error {
	err := r.q.DeleteAllFilters(ctx, membershipID)
	if err != nil {
		return NewRepoError(err, RepoErrInternal, "failed to delete all filters")
	}
	return nil
}

func (r *TagRepo) CreateTagKey(ctx context.Context, teamID int32, tagName string, dataType db.TagDataType) (db.TagKey, error) {
	tagKey, err := r.q.CreateTagKey(ctx, db.CreateTagKeyParams{
		TeamID:   teamID,
		Name:     tagName,
		DataType: dataType,
	})
	if err != nil {
		return db.TagKey{}, NewRepoError(err, RepoErrInternal, "failed to create tag key")
	}
	return tagKey, nil
}

func (r *TagRepo) CreateTagValue(ctx context.Context, tagID int32, value string) (db.TagValue, error) {
	tagValue, err := r.q.CreateTagValue(ctx, db.CreateTagValueParams{
		KeyID: tagID,
		Value: value,
	})
	if err != nil {
		return db.TagValue{}, NewRepoError(err, RepoErrInternal, "failed to create tag value")
	}
	return tagValue, nil
}

func (r *TagRepo) DeleteTagKey(ctx context.Context, tagID int32) (db.TagKey, error) {
	tagKey, err := r.q.DeleteTagKey(ctx, tagID)
	if err != nil {
		return db.TagKey{}, NewRepoError(err, RepoErrInternal, "failed to delete tag key")
	}
	return tagKey, nil
}

func (r *TagRepo) DeleteTagValue(ctx context.Context, tagValueID int32) (db.TagValue, error) {
	tagValue, err := r.q.DeleteTagValue(ctx, tagValueID)
	if err != nil {
		return db.TagValue{}, NewRepoError(err, RepoErrInternal, "failed to delete tag value")
	}
	return tagValue, nil
}

func (r *TagRepo) UpdateTagKey(ctx context.Context, tagID int32, tagName string) (db.TagKey, error) {
	tagKey, err := r.q.UpdateTagKey(ctx, db.UpdateTagKeyParams{
		ID:   tagID,
		Name: tagName,
	})
	if err != nil {
		return db.TagKey{}, NewRepoError(err, RepoErrInternal, "failed to update tag key")
	}
	return tagKey, nil
}
