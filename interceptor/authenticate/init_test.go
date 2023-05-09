package authenticate

import (
	"context"
	"github.com/justdomepaul/toolbox/services"
	"github.com/stretchr/testify/mock"
)

type testIAuthorization struct {
	mock.Mock
	services.IAuthorization
}

func (t *testIAuthorization) GetID() []byte {
	args := t.Called()
	return args.Get(0).([]byte)
}

func (t *testIAuthorization) GetClaim() interface{} {
	args := t.Called()
	return args.Get(0)
}

type testIAuthenticate struct {
	mock.Mock
}

func (t *testIAuthenticate) Authenticate(ctx context.Context, tokenFn func() (string, error), fullMethod string) (authorization services.IAuthorization, err error) {
	args := t.Called(tokenFn, fullMethod)
	return args.Get(0).(services.IAuthorization), args.Error(1)
}
