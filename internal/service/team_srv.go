package service

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skndash96/lastnight-backend/internal/db"
	"github.com/skndash96/lastnight-backend/internal/repository"
)

type TeamService struct {
	db *pgxpool.Pool
}

func NewTeamService(p *pgxpool.Pool) *TeamService {
	return &TeamService{
		db: p,
	}
}

func (s *TeamService) GetTeamsByUserID(ctx context.Context, userID int32) ([]db.GetTeamsByUserIDRow, error) {
	teamRepo := repository.NewTeamRepository(s.db)

	teams, err := teamRepo.GetTeamsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return teams, nil
}

// Do not expose this method outside of the service
// User can only join a team through JoinDefaultTeam
func (s *TeamService) joinTeam(ctx context.Context, userID, teamID int32) (*db.TeamMembership, error) {
	teamRepo := repository.NewTeamRepository(s.db)

	tm, err := teamRepo.CreateTeamMembership(ctx, userID, teamID, db.TeamUserRoleMember)

	if err != nil {
		return nil, err
	}

	return &tm, nil
}

func (s *TeamService) JoinDefaultTeam(ctx context.Context, userID int32, userEmail string) (*db.Team, error) {
	parts := strings.SplitN(strings.ToLower(userEmail), "@", 2)
	if len(parts) != 2 {
		return nil, NewSrvError(nil, SrvErrInvalidInput, "invalid email")
	}

	teamRepo := repository.NewTeamRepository(s.db)

	domain := parts[1]
	team, err := teamRepo.GetTeamByDomain(ctx, domain)

	if err != nil {
		return nil, err
	}

	_, err = s.joinTeam(ctx, userID, team.ID)
	if err != nil {
		return nil, err
	}

	return &team, nil
}
