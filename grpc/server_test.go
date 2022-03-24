package grpc

import (
	"context"
	"github.com/justdomepaul/toolbox/config"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

type testAuthenticate struct {
	mock.Mock
}

func (t *testAuthenticate) Authenticate(ctx context.Context, tokenFn func() (string, error), fullMethod string) (clientID []byte, err error) {
	args := t.Called(ctx, tokenFn, fullMethod)
	return args.Get(0).([]byte), args.Error(1)
}

type ServerSuite struct {
	suite.Suite
}

func (suite *ServerSuite) TestCreateServer() {
	suite.Equal("*grpc.Server", reflect.TypeOf(CreateServer(zap.NewExample(), config.GRPC{ALTS: true}, &testAuthenticate{})).String())
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
