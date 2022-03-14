package cloud

import (
	"cloud.google.com/go/storage"
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
	defer gostub.StubFunc(&newClient, &storage.Client{}, nil).Reset()

	option := config.Cloud{}
	suite.NoError(config.LoadFromEnv(&option))
	option.WithoutAuthentication = true
	option.GRPCInsecure = true
	_, errSession := NewSession(context.Background(), option)
	suite.NoError(errSession)
}

func (suite *ConnectServiceSuite) TestNewExtendStorageDatabase() {
	option := config.Cloud{}
	suite.NoError(config.LoadFromEnv(&option))
	result, fn, err := NewExtendStorageDatabase(
		zap.NewExample(),
		option)
	suite.NoError(err)
	defer fn()
	suite.Equal("*storage.Client", reflect.TypeOf(result).String())
	suite.Equal("func()", reflect.TypeOf(fn).String())
}

func (suite *ConnectServiceSuite) TestNewExtendStorageDatabaseNewSessionError() {
	defer gostub.StubFunc(&NewSession, &storage.Client{}, errors.New("got error")).Reset()
	option := config.Cloud{}
	suite.NoError(config.LoadFromEnv(&option))
	_, _, errExtendStorageDatabase := NewExtendStorageDatabase(
		zap.NewExample(),
		option)
	suite.Error(errExtendStorageDatabase)
}

func TestConnectServiceSuite(t *testing.T) {
	suite.Run(t, new(ConnectServiceSuite))
}
