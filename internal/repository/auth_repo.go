package repository

import (
	"context"

	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/skndash96/lastnight-backend/internal/db"
)

type AuthRepository struct {
	q *db.Queries
}

func NewAuthRepository(d db.DBTX) *AuthRepository {
	return &AuthRepository{
		q: db.New(d),
	}
}

// -------- user --------
func (r *AuthRepository) GetUserByID(ctx context.Context, id int32) (db.User, error) {
	user, err := r.q.GetUserByID(ctx, id)
	if err != nil {
		return db.User{}, NewRepoError(err, RepoErrInternal, "failed to get user by id")
	}

	return user, nil
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	user, err := r.q.GetUserByEmail(ctx, email)
	if err != nil {
		return db.User{}, NewRepoError(err, RepoErrInternal, "failed to get user by email")
	}

	return user, nil
}

func (r *AuthRepository) CreateUser(ctx context.Context, name, email string) (db.User, error) {
	user, err := r.q.CreateUser(ctx, db.CreateUserParams{
		Name: name,
		Email: email,
	})

	if err != nil {
		return db.User{}, NewRepoError(err, RepoErrInternal, "failed to create user")
	}

	return user, nil
}

// -------- account --------

func (r *AuthRepository) GetUserAccountsByID(ctx context.Context, userID int32) ([]db.Account, error) {
	accounts, err := r.q.GetUserAccountsByID(ctx, userID)
	if err != nil {
		return nil, NewRepoError(err, RepoErrInternal, "failed to get user accounts by id")
	}

	return accounts, nil
}

func (r *AuthRepository) CreateAccount(ctx context.Context, userID int32, provider, providerAccountID, password string) (db.Account, error) {
	account, err := r.q.CreateAccount(ctx, db.CreateAccountParams{
		UserID: userID,
		Provider: provider,
		ProviderAccountID: providerAccountID,
		Password: []byte(password),
	})

	if err != nil {
		return db.Account{}, NewRepoError(err, RepoErrInternal, "failed to create account")
	}

	return account, nil
}

// -------- session --------
func (r *AuthRepository) GetSessionByID(ctx context.Context, id pgtype.UUID) (db.Session, error) {
	session, err := r.q.GetSessionByID(ctx, id)
	if err != nil {
		return db.Session{}, NewRepoError(err, RepoErrInternal, "failed to get session by id")
	}

	return session, nil
}

func (r *AuthRepository) CreateSession(ctx context.Context, userID int32, email string, expiry time.Time) (*db.Session, error) {
	session, err := r.q.CreateSession(ctx, db.CreateSessionParams{
		UserID: userID,
		Email: email,
		Expiry: expiry,
	})

	if err != nil {
		return nil, NewRepoError(err, RepoErrInternal, "failed to create session")
	}

	return &session, nil
}
