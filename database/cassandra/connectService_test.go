package cassandra

import (
	"github.com/cockroachdb/errors"
	"github.com/gocql/gocql"
	"github.com/justdomepaul/toolbox/config"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"reflect"
	"strconv"
	"testing"
)

type testCassandraClusterConfig struct {
	mock.Mock
	IClusterConfig
}

func (t *testCassandraClusterConfig) CreateSession() (*gocql.Session, error) {
	args := t.Called()
	return args.Get(0).(*gocql.Session), args.Error(1)
}

type ConnectServiceSuite struct {
	suite.Suite
}

func (suite *ConnectServiceSuite) TestNewCluster() {
	suite.Equal("*gocql.ClusterConfig", reflect.TypeOf(NewCluster(&gocql.ClusterConfig{})).String())
}

func (suite *ConnectServiceSuite) TestCreateSession() {
	testCassandraSession := &gocql.Session{}
	testCassandraClusterConfig := &testCassandraClusterConfig{}
	testCassandraClusterConfig.On("CreateSession").Return(testCassandraSession, nil)

	defer gostub.StubFunc(&extendOptions, nil).Reset()
	defer gostub.StubFunc(&NewCluster, testCassandraClusterConfig).Reset()
	session, err := NewSession(config.Cassandra{
		CassandraHosts: []string{"localhost:8080"},
	})
	suite.NoError(err)
	suite.Equal("*gocql.Session", reflect.TypeOf(session).String())
}

func (suite *ConnectServiceSuite) TestCreateSessionExtendOptionsError() {
	defer gostub.StubFunc(&extendOptions, errors.New("got error")).Reset()
	_, err := NewSession(config.Cassandra{
		CassandraHosts: []string{"localhost:8080"},
	})
	suite.Error(err)
}

func (suite *ConnectServiceSuite) TestExtendOptions() {
	mockStrPort := "9090"
	mockPort, errPort := strconv.Atoi("9090")
	suite.NoError(errPort)
	cg := &gocql.ClusterConfig{}
	mockSpecification := config.Cassandra{
		CassandraKeySpace:          "testKeySpace",
		CassandraPort:              mockStrPort,
		CassandraHosts:             []string{"localhost"},
		CassandraForcePassword:     true,
		CassandraUsername:          "testUserName",
		CassandraPassword:          "testPassword",
		CassandraTimeout:           0,
		CassandraConnectionTimeout: 0,
		ReconnectInterval:          0,
		CassandraNumConnection:     0,
		DisableInitialHostLookup:   false,
	}
	suite.NoError(extendOptions(cg, mockSpecification))
	suite.Equal(mockSpecification.CassandraKeySpace, cg.Keyspace)
	suite.Equal(mockPort, cg.Port)
	suite.Equal(mockSpecification.CassandraUsername, cg.Authenticator.(gocql.PasswordAuthenticator).Username)
	suite.Equal(mockSpecification.CassandraPassword, cg.Authenticator.(gocql.PasswordAuthenticator).Password)
	suite.Equal(mockSpecification.CassandraConnectionTimeout, cg.Timeout)
	suite.Equal(mockSpecification.ReconnectInterval, cg.ReconnectInterval)
	suite.Equal(mockSpecification.CassandraConnectionTimeout, cg.ConnectTimeout)
	suite.Equal(mockSpecification.CassandraNumConnection, cg.NumConns)
	suite.Equal(mockSpecification.DisableInitialHostLookup, cg.DisableInitialHostLookup)
}

func (suite *ConnectServiceSuite) TestExtendOptionsPortError() {
	cg := &gocql.ClusterConfig{}
	mockSpecification := config.Cassandra{
		CassandraKeySpace:          "testKeySpace",
		CassandraPort:              "testPort",
		CassandraHosts:             []string{"localhost"},
		CassandraForcePassword:     true,
		CassandraUsername:          "testUserName",
		CassandraPassword:          "testPassword",
		CassandraTimeout:           0,
		CassandraConnectionTimeout: 0,
		ReconnectInterval:          0,
		CassandraNumConnection:     0,
		DisableInitialHostLookup:   false,
	}
	suite.Error(extendOptions(cg, mockSpecification))
}

func (suite *ConnectServiceSuite) TestNewExtendCockroachDatabase() {
	testCassandraSession := &gocql.Session{}
	testCassandraClusterConfig := &testCassandraClusterConfig{}
	testCassandraClusterConfig.On("CreateSession").Return(testCassandraSession, nil)

	defer gostub.StubFunc(&extendOptions, nil).Reset()
	defer gostub.StubFunc(&NewCluster, testCassandraClusterConfig).Reset()

	result, fn, err := NewExtendCassandraDatabase(
		zap.NewExample(),
		config.Cassandra{
			CassandraHosts: []string{"test", "test2", "test3"},
		})
	suite.NoError(err)
	suite.Equal("*gocql.Session", reflect.TypeOf(result).String())
	suite.Equal("func()", reflect.TypeOf(fn).String())
	suite.NotPanics(func() {
		fn()
	})
}

func (suite *ConnectServiceSuite) TestNewExtendCockroachDatabaseNewSessionError() {
	defer gostub.StubFunc(&NewSession, nil, errors.New("got error")).Reset()

	_, _, err := NewExtendCassandraDatabase(
		zap.NewExample(),
		config.Cassandra{
			CassandraHosts: []string{"test", "test2", "test3"},
		})
	suite.Error(err)
}

func TestConnectServiceSuite(t *testing.T) {
	suite.Run(t, new(ConnectServiceSuite))
}
