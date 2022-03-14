package config

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

type CassandraSuite struct {
	suite.Suite
	CassandraKeySpace          string
	CassandraPort              string
	CassandraHosts             []string
	CassandraForcePassword     bool
	CassandraUsername          string
	CassandraPassword          string
	CassandraTimeout           time.Duration
	CassandraConnectionTimeout time.Duration
	ReconnectInterval          time.Duration
	CassandraNumConnection     int
	DisableInitialHostLookup   bool
}

func (suite *CassandraSuite) SetupSuite() {
	os.Clearenv()
	suite.CassandraKeySpace = "testCassandraKeySpace"
	suite.CassandraPort = "testCassandraPort"
	suite.CassandraHosts = []string{"testCassandraHosts"}
	suite.CassandraForcePassword = true
	suite.CassandraUsername = "testCassandraUsername"
	suite.CassandraPassword = "testCassandraPassword"
	suite.CassandraTimeout = 10 * time.Second
	suite.CassandraConnectionTimeout = 10 * time.Second
	suite.ReconnectInterval = 10 * time.Second
	suite.CassandraNumConnection = 100
	suite.DisableInitialHostLookup = true

	suite.NoError(os.Setenv("CASSANDRA_KEY_SPACE", suite.CassandraKeySpace))
	suite.NoError(os.Setenv("CASSANDRA_PORT", suite.CassandraPort))
	suite.NoError(os.Setenv("CASSANDRA_HOSTS", strings.Join(suite.CassandraHosts, ",")))
	suite.NoError(os.Setenv("CASSANDRA_FORCE_PASSWORD", strconv.FormatBool(suite.CassandraForcePassword)))
	suite.NoError(os.Setenv("CASSANDRA_USERNAME", suite.CassandraUsername))
	suite.NoError(os.Setenv("CASSANDRA_PASSWORD", suite.CassandraPassword))
	suite.NoError(os.Setenv("CASSANDRA_TIMEOUT", fmt.Sprint(suite.CassandraTimeout)))
	suite.NoError(os.Setenv("CASSANDRA_CONNECTION_TIMEOUT", fmt.Sprint(suite.CassandraConnectionTimeout)))
	suite.NoError(os.Setenv("RECONNECT_INTERVAL", fmt.Sprint(suite.ReconnectInterval)))
	suite.NoError(os.Setenv("CASSANDRA_NUM_CONNECTION", strconv.FormatInt(int64(suite.CassandraNumConnection), 10)))
	suite.NoError(os.Setenv("DISABLE_INITIAL_HOST_LOOKUP", strconv.FormatBool(suite.DisableInitialHostLookup)))

}

func (suite *CassandraSuite) TestDefaultOption() {
	cassandra := &Cassandra{}
	suite.NoError(LoadFromEnv(cassandra))
	suite.Equal(suite.CassandraKeySpace, cassandra.CassandraKeySpace)
	suite.Equal(suite.CassandraPort, cassandra.CassandraPort)
	suite.Equal(suite.CassandraHosts, cassandra.CassandraHosts)
	suite.Equal(suite.CassandraForcePassword, cassandra.CassandraForcePassword)
	suite.Equal(suite.CassandraUsername, cassandra.CassandraUsername)
	suite.Equal(suite.CassandraPassword, cassandra.CassandraPassword)
	suite.Equal(suite.CassandraTimeout, cassandra.CassandraTimeout)
	suite.Equal(suite.CassandraConnectionTimeout, cassandra.CassandraConnectionTimeout)
	suite.Equal(suite.ReconnectInterval, cassandra.ReconnectInterval)
	suite.Equal(suite.CassandraNumConnection, cassandra.CassandraNumConnection)
	suite.Equal(suite.DisableInitialHostLookup, cassandra.DisableInitialHostLookup)
}

func TestCassandraSuite(t *testing.T) {
	suite.Run(t, new(CassandraSuite))
}
