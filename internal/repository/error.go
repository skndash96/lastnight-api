package repository

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

type RepoError struct {
	internal error
	Kind     RepoErrKind
	Message  string
}

type RepoErrKind string

const (
	// Client/input related issues detected at the DB layer (constraint, validation)
	RepoErrInvalidInput RepoErrKind = "invalid_input"

	// Conflicts or constraint violations (unique, FK, etc.)
	RepoErrConflict RepoErrKind = "conflict"

	// Unexpected database/system failure
	RepoErrInternal RepoErrKind = "internal_error"
)

func (e *RepoError) Error() string {
	return e.Message
}

func (e *RepoError) Unwrap() error {
	return e.internal
}

// NewRepoError constructs a RepoError. If the provided error is a Postgres error,
// it will be mapped to an appropriate RepoErrKind and Message.
func NewRepoError(err error, kind RepoErrKind, message string) *RepoError {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return mapPgErrToRepoError(pgErr)
	}

	return &RepoError{
		internal: err,
		Kind:     kind,
		Message:  message,
	}
}

// mapPgErrToRepoError maps common Postgres error codes to repository error kinds/messages.
// Reference: https://www.postgresql.org/docs/current/errcodes-appendix.html
func mapPgErrToRepoError(err *pgconn.PgError) *RepoError {
	switch err.Code {
	// Integrity constraint violations
	case "23505": // unique_violation
		return &RepoError{
			internal: err,
			Kind:     RepoErrConflict,
			Message:  "conflicting values",
		}
	case "23503": // foreign_key_violation
		return &RepoError{
			internal: err,
			Kind:     RepoErrConflict,
			Message:  "conflicting relation values",
		}
	case "23502": // not_null_violation
		return &RepoError{
			internal: err,
			Kind:     RepoErrInvalidInput,
			Message:  "missing required values",
		}
	case "22001": // string_data_right_truncation
		return &RepoError{
			internal: err,
			Kind:     RepoErrInvalidInput,
			Message:  "input exceeds allowed length",
		}
	case "22P02": // invalid_text_representation (e.g., invalid UUID)
		return &RepoError{
			internal: err,
			Kind:     RepoErrInvalidInput,
			Message:  "invalid input representation",
		}

	// Syntax / undefined objects typically indicate internal dev issues
	case "42601": // syntax_error
	case "42703": // undefined_column
	case "42P01": // undefined_table
		return &RepoError{
			internal: err,
			Kind:     RepoErrInternal,
			Message:  "database query error",
		}
	}

	return &RepoError{
		internal: err,
		Kind:     RepoErrInternal,
		Message:  "unknown database error",
	}
}
