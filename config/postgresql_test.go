package config

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

type PostgresqlSuite struct {
	suite.Suite
	PostgresqlUsername  string
	PostgresqlPassword  string
	PostgresqlHost      string
	PostgresqlPort      string
	PostgresqlDatabase  string
	PostgresqlSSLMode   string
	PostgresqlURL       string
	PostgresqlTxTimeout time.Duration
}

func (suite *PostgresqlSuite) SetupSuite() {
	os.Clearenv()
	suite.PostgresqlUsername = "testPostgresqlUsername"
	suite.PostgresqlPassword = "testPostgresqlPassword"
	suite.PostgresqlHost = "testPostgresqlHost"
	suite.PostgresqlPort = "testPostgresqlPort"
	suite.PostgresqlDatabase = "testPostgresqlDatabase"
	suite.PostgresqlSSLMode = "testPostgresqlSSLMode"
	suite.PostgresqlURL = "testPostgresqlURL"
	suite.PostgresqlTxTimeout = 20 * time.Second
	suite.NoError(os.Setenv("POSTGRESQL_USERNAME", suite.PostgresqlUsername))
	suite.NoError(os.Setenv("POSTGRESQL_PASSWORD", suite.PostgresqlPassword))
	suite.NoError(os.Setenv("POSTGRESQL_HOST", suite.PostgresqlHost))
	suite.NoError(os.Setenv("POSTGRESQL_PORT", suite.PostgresqlPort))
	suite.NoError(os.Setenv("POSTGRESQL_DATABASE", suite.PostgresqlDatabase))
	suite.NoError(os.Setenv("POSTGRESQL_SSL_MODE", suite.PostgresqlSSLMode))
	suite.NoError(os.Setenv("POSTGRESQL_URL", suite.PostgresqlURL))
	suite.NoError(os.Setenv("POSTGRESQL_TX_TIMEOUT", fmt.Sprint(suite.PostgresqlTxTimeout)))
}

func (suite *PostgresqlSuite) TestDefaultOption() {
	postgresql := &Postgresql{}
	suite.NoError(LoadFromEnv(postgresql))
	suite.Equal(suite.PostgresqlUsername, postgresql.PostgresqlUsername)
	suite.Equal(suite.PostgresqlPassword, postgresql.PostgresqlPassword)
	suite.Equal(suite.PostgresqlHost, postgresql.PostgresqlHost)
	suite.Equal(suite.PostgresqlPort, postgresql.PostgresqlPort)
	suite.Equal(suite.PostgresqlDatabase, postgresql.PostgresqlDatabase)
	suite.Equal(suite.PostgresqlSSLMode, postgresql.PostgresqlSSLMode)
	suite.Equal(suite.PostgresqlURL, postgresql.PostgresqlURL)
	suite.Equal(suite.PostgresqlTxTimeout, postgresql.PostgresqlTxTimeout)
}

func TestPostgresqlSuite(t *testing.T) {
	suite.Run(t, new(PostgresqlSuite))
}
