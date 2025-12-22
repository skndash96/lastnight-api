package service

type SrvError struct {
	internal error
	Kind     SrvErrKind
	Message  string
}

type SrvErrKind string

const (
	// Client fault (input or business logic)
	SrvErrInvalidInput SrvErrKind = "invalid_input" // malformed data, violated constraints
	SrvErrUnauthorized SrvErrKind = "unauthorized"  // auth failed
	SrvErrForbidden    SrvErrKind = "forbidden"     // auth ok, but access denied
	SrvErrNotFound     SrvErrKind = "not_found"     // lookup failed
	SrvErrConflict     SrvErrKind = "conflict"      // duplicate resource / violation (email exists)

	// Non-client, unexpected runtime/system errors
	SrvErrInternal SrvErrKind = "internal_error" // unexpected application failure
)

func (e *SrvError) Error() string {
	return e.Message
}

func (e *SrvError) Unwrap() error {
	return e.internal
}

func NewSrvError(err error, kind SrvErrKind, message string) *SrvError {
	return &SrvError{
		internal: err,
		Kind:     kind,
		Message:  message,
	}
}
