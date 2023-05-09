package services

import (
	"context"
)

type IAuthorization interface {
	GetID() []byte
	// GetClaim response must assert to specify jwt claim struct
	GetClaim() interface{}
}

// IAuthenticate interface
// returns error, ErrInWhitelist, ErrInvalidArguments, ErrUnauthenticated, ErrDeny, ErrNoRefreshToken, ErrScopeNotExist, ErrOutOfScopes
// if success return id
type IAuthenticate interface {
	Authenticate(ctx context.Context, tokenFn func() (string, error), fullMethod string) (authorization IAuthorization, err error)
}
