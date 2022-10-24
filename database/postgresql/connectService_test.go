package postgresql

import (
	"github.com/cockroachdb/errors"
	"github.com/jmoiron/sqlx"
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

func (suite *ConnectServiceSuite) TestNewSessionURL() {
	defer gostub.StubFunc(&Connect, &sqlx.DB{}, nil).Reset()

	_, err := NewSession(config.Postgresql{
		PostgresqlURL: "testURL",
	})
	suite.NoError(err)
}

func (suite *ConnectServiceSuite) TestNewSessionParameters() {
	defer gostub.StubFunc(&Connect, &sqlx.DB{}, nil).Reset()

	_, err := NewSession(config.Postgresql{
		PostgresqlUsername: "testPostgresqlUsername",
		PostgresqlPassword: "testPostgresqlPassword",
		PostgresqlHost:     "testPostgresqlHost",
		PostgresqlPort:     "testPostgresqlPort",
		PostgresqlDatabase: "testPostgresqlDatabase",
		PostgresqlSSLMode:  "testPostgresqlSSLMode",
	})
	suite.NoError(err)
}

func (suite *ConnectServiceSuite) TestNewExtendPostgresqlDatabase() {
	defer gostub.StubFunc(&NewSession, new(testISession), nil).Reset()
	result, fn, err := NewExtendPostgresqlDatabase(
		zap.NewExample(),
		config.Postgresql{})
	suite.NoError(err)
	suite.Equal("*postgresql.testISession", reflect.TypeOf(result).String())
	suite.Equal("func()", reflect.TypeOf(fn).String())
}

func (suite *ConnectServiceSuite) TestNewExtendPostgresqlDatabaseNewSessionError() {
	defer gostub.StubFunc(&NewSession, nil, errors.New("got error")).Reset()
	_, _, err := NewExtendPostgresqlDatabase(
		zap.NewExample(),
		config.Postgresql{})
	suite.Error(err)
}

func (suite *ConnectServiceSuite) TestNewExtendPostgresqlDatabaseSessionCloseError() {
	testISession := &testISession{}
	testISession.On("Close").Return(errors.New("got error"))

	defer gostub.StubFunc(&NewSession, testISession, nil).Reset()

	_, fn, err := NewExtendPostgresqlDatabase(zap.NewExample(), config.Postgresql{})
	suite.NoError(err)
	suite.NotPanics(func() {
		fn()
	})
}

func TestConnectServiceSuite(t *testing.T) {
	suite.Run(t, new(ConnectServiceSuite))
}
