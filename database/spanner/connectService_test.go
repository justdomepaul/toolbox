package spanner

import (
	"cloud.google.com/go/spanner"
	"context"
	"github.com/cockroachdb/errors"
	"github.com/justdomepaul/toolbox/config"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
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
	t := suite.T()
	defer gostub.StubFunc(&newClient, &spanner.Client{}, nil).Reset()

	option := config.Spanner{}
	suite.NoError(config.LoadFromEnv(&option))
	option.WithoutAuthentication = true
	option.GRPCInsecure = true
	_, errSession := NewSession(context.Background(), option)
	assert.NoError(t, errSession)
}

func (suite *ConnectServiceSuite) TestNewExtendSpannerDatabase() {
	t := suite.T()
	defer gostub.StubFunc(&NewSession, new(testISession), nil).Reset()
	option := config.Spanner{}
	suite.NoError(config.LoadFromEnv(&option))
	result, fn, err := NewExtendSpannerDatabase(
		zap.NewExample(),
		option)
	assert.NoError(t, err)
	defer fn()
	assert.Equal(t, "*spanner.testISession", reflect.TypeOf(result).String())
	assert.Equal(t, "func()", reflect.TypeOf(fn).String())
}

func (suite *ConnectServiceSuite) TestNewExtendSpannerDatabaseNewSessionError() {
	t := suite.T()
	defer gostub.StubFunc(&NewSession, new(testISession), errors.New("got error")).Reset()
	option := config.Spanner{}
	suite.NoError(config.LoadFromEnv(&option))
	_, _, errExtendSpannerDatabase := NewExtendSpannerDatabase(
		zap.NewExample(),
		option)
	assert.Error(t, errExtendSpannerDatabase)
}

func TestConnectServiceSuite(t *testing.T) {
	suite.Run(t, new(ConnectServiceSuite))
}
