package bunt

import (
	"github.com/cockroachdb/errors"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
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

func (suite *ConnectServiceSuite) TestCreateSession() {
	t := suite.T()
	session, err := NewSession(":memory:")
	assert.NoError(t, err)
	assert.Equal(t, "*buntdb.DB", reflect.TypeOf(session).String())
}

func (suite *ConnectServiceSuite) TestNewExtendBuntDatabase() {
	t := suite.T()

	defer gostub.StubFunc(&NewSession, new(testISession), nil).Reset()

	result, fn, err := NewExtendBuntDatabase(zap.NewExample())
	assert.NoError(t, err)
	assert.Equal(t, "*bunt.testISession", reflect.TypeOf(result).String())
	assert.Equal(t, "func()", reflect.TypeOf(fn).String())
	assert.NotPanics(t, func() {
		fn()
	})
}

func (suite *ConnectServiceSuite) TestNewExtendBuntDatabaseNewSessionError() {
	t := suite.T()

	defer gostub.StubFunc(&NewSession, new(testISession), errors.New("got error")).Reset()

	_, _, err := NewExtendBuntDatabase(zap.NewExample())
	assert.Error(t, err)
}

func (suite *ConnectServiceSuite) TestNewExtendBuntDatabaseSessionCloseError() {
	t := suite.T()

	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	logger := zap.New(observedZapCore, zap.Fields(zap.String("system", "Mock system")))

	testISession := &testISession{}
	testISession.On("Close").Return(errors.New("got error"))

	defer gostub.StubFunc(&NewSession, testISession, nil).Reset()

	_, fn, err := NewExtendBuntDatabase(logger)
	assert.NoError(t, err)
	fn()
	require.Equal(suite.T(), 1, observedLogs.Len())
	firstLog := observedLogs.All()[0]
	suite.Equal("Bunt init complete", firstLog.Message)
}

func TestConnectServiceSuite(t *testing.T) {
	suite.Run(t, new(ConnectServiceSuite))
}
