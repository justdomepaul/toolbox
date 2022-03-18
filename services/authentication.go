package services

import (
	"context"
)

// Authenticate interface
// Authenticate returns
// ErrInWhitelist, ErrInvalidArguments, ErrUnauthenticated, ErrDeny, ErrNoRefreshToken, ErrScopeNotExist, ErrOutOfScopes
// if success return id
type Authenticate interface {
	Authenticate(ctx context.Context, tokenFn func() (string, error), fullMethod string) (clientID []byte, err error)
}
