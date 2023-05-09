package stateful

import (
	"context"
	"fmt"
	"github.com/justdomepaul/toolbox/array"
	"github.com/justdomepaul/toolbox/config"
	"github.com/justdomepaul/toolbox/errorhandler"
	"github.com/justdomepaul/toolbox/jwt"
	"github.com/justdomepaul/toolbox/services"
	"strings"
)

var (
	newToken = jwt.NewCommon
)

func NewAuthorization(id []byte, claim *jwt.Common) *Authorization {
	return &Authorization{
		ID:    id,
		Claim: claim,
	}
}

type Authorization struct {
	ID    []byte
	Claim *jwt.Common
}

func (a *Authorization) GetID() []byte {
	return a.ID
}

func (a *Authorization) GetClaim() interface{} {
	return a.Claim
}

func NewAuthentication(gRPC config.GRPC, jwt jwt.IJWT) (*Authentication, error) {
	return &Authentication{
		allowedList: gRPC.AllowedList,
		j:           jwt,
	}, nil
}

type Authentication struct {
	allowedList []string
	j           jwt.IJWT
}

func (s *Authentication) Authenticate(ctx context.Context, tokenFn func() (string, error), fullMethod string) (authorization services.IAuthorization, err error) {
	for _, term := range s.allowedList {
		if strings.HasPrefix(fullMethod, term) {
			return nil, errorhandler.ErrInWhitelist
		}
	}
	token, err := tokenFn()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", errorhandler.ErrUnauthenticated, err.Error())
	}

	claim := newToken(jwt.NewClaimsBuilder().Build())
	if err := s.j.VerifyToken(token, claim); err != nil {
		return nil, fmt.Errorf("%w: %s", errorhandler.ErrUnauthenticated, err.Error())
	}

	if _, exist := array.Find(claim.Scopes, fullMethod); !exist {
		return nil, errorhandler.ErrOutOfScopes
	}

	return NewAuthorization(claim.ClientID, claim), nil
}
