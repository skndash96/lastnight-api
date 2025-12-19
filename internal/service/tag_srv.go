package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skndash96/lastnight-backend/internal/db"
	"github.com/skndash96/lastnight-backend/internal/helpers"
	"github.com/skndash96/lastnight-backend/internal/repository"
)

type TagService struct {
	db *pgxpool.Pool
}

func NewTagService(p *pgxpool.Pool) *TagService {
	return &TagService{
		db: p,
	}
}

func (s *TagService) ListFilters(ctx context.Context, membershipID int32) ([]db.Tag, error) {
	tagRepo := repository.NewTagRepo(s.db)
	tags, err := tagRepo.ListFilters(ctx, membershipID)
	if err != nil {
		return nil, NewSrvError(err, SrvErrInternal, "Failed to list tags")
	}
	return tags, nil
}

func (s *TagService) CreateTagKey(ctx context.Context, teamID int32, name string, dataType db.TagDataType) (*db.TagKey, error) {
	tagRepo := repository.NewTagRepo(s.db)
	tag, err := tagRepo.CreateTagKey(ctx, teamID, name, dataType)
	if err != nil {
		if helpers.IsUniqueViolation(err) {
			return nil, NewSrvError(err, SrvErrInvalidInput, "tag value already exists")
		}
		return nil, NewSrvError(err, SrvErrInternal, "failed to create tag key")
	}
	return &tag, nil
}

func (s *TagService) UpdateTagKey(ctx context.Context, tagID int32, name string) (*db.TagKey, error) {
	tagRepo := repository.NewTagRepo(s.db)
	tag, err := tagRepo.UpdateTagKey(ctx, tagID, name)
	if err != nil {
		if helpers.IsUniqueViolation(err) {
			return nil, NewSrvError(err, SrvErrInvalidInput, "tag value already exists")
		}
		return nil, NewSrvError(err, SrvErrInternal, "failed to update tag key")
	}
	return &tag, nil
}

func (s *TagService) DeleteTagKey(ctx context.Context, tagID int32) (*db.TagKey, error) {
	tagRepo := repository.NewTagRepo(s.db)
	tag, err := tagRepo.DeleteTagKey(ctx, tagID)
	if err != nil {
		return nil, NewSrvError(err, SrvErrInternal, "failed to delete tag key")
	}
	return &tag, nil
}

func (s *TagService) CreateTagValue(ctx context.Context, tagID int32, value string) (*db.TagValue, error) {
	tagValueRepo := repository.NewTagRepo(s.db)
	tagValue, err := tagValueRepo.CreateTagValue(ctx, tagID, value)
	if err != nil {
		if helpers.IsUniqueViolation(err) {
			return nil, NewSrvError(err, SrvErrInvalidInput, "tag value already exists")
		}
		return nil, NewSrvError(err, SrvErrInternal, "failed to create tag value")
	}
	return &tagValue, nil
}

func (s *TagService) DeleteTagValue(ctx context.Context, tagValueID int32) (*db.TagValue, error) {
	tagValueRepo := repository.NewTagRepo(s.db)
	tagValue, err := tagValueRepo.DeleteTagValue(ctx, tagValueID)
	if err != nil {
		return nil, NewSrvError(err, SrvErrInternal, "failed to delete tag value")
	}
	return &tagValue, nil
}
