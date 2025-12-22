package repository

import (
	"context"

	"github.com/skndash96/lastnight-backend/internal/db"
)

type TeamRepository struct {
	q *db.Queries
}

func NewTeamRepository(d db.DBTX) *TeamRepository {
	return &TeamRepository{
		q: db.New(d),
	}
}

func (r *TeamRepository) GetTeamsByUserID(ctx context.Context, id int32) ([]db.GetTeamsByUserIDRow, error) {
	teams, err := r.q.GetTeamsByUserID(ctx, id)
	if err != nil {
		return nil, NewRepoError(err, RepoErrInternal, "failed to get teams by user ID")
	}
	return teams, nil
}

func (r *TeamRepository) GetTeamByDomain(ctx context.Context, d string) (db.Team, error) {
	team, err := r.q.GetTeamByDomain(ctx, d)
	if err != nil {
		return db.Team{}, NewRepoError(err, RepoErrInternal, "failed to get team by domain")
	}
	return team, nil
}

func (r *TeamRepository) GetTeamMembershipByUserID(ctx context.Context, user_id, team_id int32) (db.TeamMembership, error) {
	membership, err := r.q.GetTeamMembershipByUserID(ctx, db.GetTeamMembershipByUserIDParams{
		UserID: user_id,
		TeamID: team_id,
	})
	if err != nil {
		return db.TeamMembership{}, NewRepoError(err, RepoErrInternal, "failed to get team membership by user ID")
	}
	return membership, nil
}

func (r *TeamRepository) CreateTeamMembership(ctx context.Context, user_id, team_id int32, role db.TeamUserRole) (db.TeamMembership, error) {
	membership, err := r.q.CreateTeamMembership(ctx, db.CreateTeamMembershipParams{
		UserID: user_id,
		TeamID: team_id,
		Role:   role,
	})
	if err != nil {
		return db.TeamMembership{}, NewRepoError(err, RepoErrInternal, "failed to create team membership")
	}
	return membership, nil
}
