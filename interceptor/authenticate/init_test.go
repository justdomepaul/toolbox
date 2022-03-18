package authenticate

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type testIAuthenticate struct {
	mock.Mock
}

func (t *testIAuthenticate) Authenticate(ctx context.Context, tokenFn func() (string, error), fullMethod string) (clientID []byte, err error) {
	args := t.Called(tokenFn, fullMethod)
	return args.Get(0).([]byte), args.Error(1)
}
