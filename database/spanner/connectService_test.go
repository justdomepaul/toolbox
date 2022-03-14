package spanner

import (
	"cloud.google.com/go/spanner"
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

func (t *testISession) Close() {
	t.Called()
}

type ConnectServiceSuite struct {
	suite.Suite
}

func (suite *ConnectServiceSuite) TestNewSession() {
	defer gostub.StubFunc(&newClient, &spanner.Client{}, nil).Reset()

	option := config.Spanner{}
	suite.NoError(config.LoadFromEnv(&option))
	option.WithoutAuthentication = true
	option.GRPCInsecure = true
	_, errSession := NewSession(context.Background(), option)
	suite.NoError(errSession)
}

func (suite *ConnectServiceSuite) TestNewExtendSpannerDatabase() {
	defer gostub.StubFunc(&NewSession, new(testISession), nil).Reset()
	option := config.Spanner{}
	suite.NoError(config.LoadFromEnv(&option))
	result, fn, err := NewExtendSpannerDatabase(
		zap.NewExample(),
		option)
	suite.NoError(err)
	defer fn()
	suite.Equal("*spanner.testISession", reflect.TypeOf(result).String())
	suite.Equal("func()", reflect.TypeOf(fn).String())
}

func (suite *ConnectServiceSuite) TestNewExtendSpannerDatabaseNewSessionError() {
	defer gostub.StubFunc(&NewSession, new(testISession), errors.New("got error")).Reset()
	option := config.Spanner{}
	suite.NoError(config.LoadFromEnv(&option))
	_, _, errExtendSpannerDatabase := NewExtendSpannerDatabase(
		zap.NewExample(),
		option)
	suite.Error(errExtendSpannerDatabase)
}

func TestConnectServiceSuite(t *testing.T) {
	suite.Run(t, new(ConnectServiceSuite))
}
