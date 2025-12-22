package provider

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/skndash96/lastnight-backend/internal/auth"
	"github.com/skndash96/lastnight-backend/internal/config"
	"github.com/skndash96/lastnight-backend/internal/repository"
)

type SessionProvider struct {
	sCfg     config.SessionConfig
	authRepo *repository.AuthRepository
}

func NewSessionProvider(sCfg config.SessionConfig, authRepo *repository.AuthRepository) auth.TokenProvider {
	return &SessionProvider{
		sCfg:     sCfg,
		authRepo: authRepo,
	}
}

func (p *SessionProvider) GenerateToken(ctx context.Context, userID int32, email string) (string, error) {
	session, err := p.authRepo.CreateSession(ctx, userID, email, time.Now().Add(p.sCfg.Expiry))
	if err != nil {
		return "", err
	}
	return session.ID.String(), nil
}

func (p *SessionProvider) ValidateToken(ctx context.Context, token string) (*auth.Actor, error) {
	sessionID := pgtype.UUID{}
	actor := auth.Actor{}

	err := sessionID.Scan(token)
	if err != nil {
		return &actor, err
	}

	session, err := p.authRepo.GetSessionByID(ctx, sessionID)
	if err != nil {
		return &actor, err
	}

	actor.UserID = session.UserID
	actor.Email = session.Email

	return &actor, nil
}
