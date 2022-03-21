package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/cockroachdb/errors"
	"github.com/justdomepaul/toolbox/config"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

type testISession struct {
	mock.Mock
	ISession
}

func (t *testISession) Close() error {
	args := t.Called()
	return args.Error(0)
}

type ConnectServiceSuite struct {
	suite.Suite
}

func (suite *ConnectServiceSuite) TestNewSession() {
	defer gostub.StubFunc(&newClient, &pubsub.Client{}, nil).Reset()

	option := config.PubSub{}
	suite.NoError(config.LoadFromEnv(&option))
	option.WithoutAuthentication = true
	option.GRPCInsecure = true
	_, errSession := NewSession(context.Background(), option)
	suite.NoError(errSession)
}

func (suite *ConnectServiceSuite) TestNewExtendPubSubDatabase() {
	defer gostub.StubFunc(&NewSession, new(testISession), nil).Reset()
	option := config.PubSub{}
	suite.NoError(config.LoadFromEnv(&option))
	result, fn, err := NewExtendPubSubDatabase(
		zap.NewExample(),
		option)
	suite.NoError(err)
	defer fn()
	suite.Equal("*pubsub.testISession", reflect.TypeOf(result).String())
	suite.Equal("func()", reflect.TypeOf(fn).String())
}

func (suite *ConnectServiceSuite) TestNewExtendPubSubDatabaseNewSessionError() {
	defer gostub.StubFunc(&NewSession, new(testISession), errors.New("got error")).Reset()
	option := config.PubSub{}
	suite.NoError(config.LoadFromEnv(&option))
	_, _, errExtendPubSubDatabase := NewExtendPubSubDatabase(
		zap.NewExample(),
		option)
	suite.Error(errExtendPubSubDatabase)
}

func TestConnectServiceSuite(t *testing.T) {
	suite.Run(t, new(ConnectServiceSuite))
}
