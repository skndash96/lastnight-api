package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/skndash96/lastnight-backend/internal/auth"
	"github.com/skndash96/lastnight-backend/internal/db"
	"github.com/skndash96/lastnight-backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Register(ctx context.Context, name, email, password string) (string, error)
}

type authService struct {
	db            *pgxpool.Pool
	tokenProvider auth.TokenProvider
}

func NewAuthService(db *pgxpool.Pool, tokenProvider auth.TokenProvider) AuthService {
	return &authService{
		db:            db,
		tokenProvider: tokenProvider,
	}
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	authRepo := repository.NewAuthRepository(s.db)

	u, err := authRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	accounts, err := authRepo.GetUserAccountsByID(ctx, u.ID)
	if err != nil {
		return "", err
	}

	var acc *db.Account
	for _, a := range accounts {
		if a.Provider == "local" {
			acc = &a
			break
		}
	}

	if acc == nil {
		return "", NewSrvError(err, SrvErrInvalidInput, "local account not found")
	}

	if err := bcrypt.CompareHashAndPassword(acc.Password, []byte(password)); err != nil {
		return "", NewSrvError(err, SrvErrInvalidInput, "invalid credentials")
	}

	token, err := s.tokenProvider.GenerateToken(ctx, acc.UserID, u.Email)
	if err != nil {
		return "", NewSrvError(err, SrvErrInternal, "failed to generate token")
	}

	return token, nil
}

func (s *authService) Register(ctx context.Context, name, email, password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", NewSrvError(err, SrvErrInternal, "failed to hash password")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return "", NewSrvError(err, SrvErrInternal, "failed to start transaction")
	}
	defer tx.Rollback(ctx)

	authRepo := repository.NewAuthRepository(tx)

	user, err := authRepo.CreateUser(ctx, name, email)
	if err != nil {
		return "", err
	}

	acc, err := authRepo.CreateAccount(ctx, user.ID, "local", email, string(passwordHash))
	if err != nil {
		return "", err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", NewSrvError(err, SrvErrInternal, "failed to commit transaction")
	}

	token, err := s.tokenProvider.GenerateToken(ctx, acc.UserID, user.Email)
	if err != nil {
		return "", NewSrvError(err, SrvErrInternal, "failed to generate token")
	}

	return token, nil
}
