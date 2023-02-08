package config

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

type PostgresSuite struct {
	suite.Suite
	PostgresUsername  string
	PostgresPassword  string
	PostgresHost      string
	PostgresPort      string
	PostgresDatabase  string
	PostgresSSLMode   string
	PostgresURL       string
	PostgresTxTimeout time.Duration
}

func (suite *PostgresSuite) SetupSuite() {
	os.Clearenv()
	suite.PostgresUsername = "testPostgresUsername"
	suite.PostgresPassword = "testPostgresPassword"
	suite.PostgresHost = "testPostgresHost"
	suite.PostgresPort = "testPostgresPort"
	suite.PostgresDatabase = "testPostgresDatabase"
	suite.PostgresSSLMode = "testPostgresSSLMode"
	suite.PostgresURL = "testPostgresURL"
	suite.PostgresTxTimeout = 20 * time.Second
	suite.NoError(os.Setenv("POSTGRES_USERNAME", suite.PostgresUsername))
	suite.NoError(os.Setenv("POSTGRES_PASSWORD", suite.PostgresPassword))
	suite.NoError(os.Setenv("POSTGRES_HOST", suite.PostgresHost))
	suite.NoError(os.Setenv("POSTGRES_PORT", suite.PostgresPort))
	suite.NoError(os.Setenv("POSTGRES_DATABASE", suite.PostgresDatabase))
	suite.NoError(os.Setenv("POSTGRES_SSL_MODE", suite.PostgresSSLMode))
	suite.NoError(os.Setenv("POSTGRES_URL", suite.PostgresURL))
	suite.NoError(os.Setenv("POSTGRES_TX_TIMEOUT", fmt.Sprint(suite.PostgresTxTimeout)))
}

func (suite *PostgresSuite) TestDefaultOption() {
	postgres := &Postgres{}
	suite.NoError(LoadFromEnv(postgres))
	suite.Equal(suite.PostgresUsername, postgres.PostgresUsername)
	suite.Equal(suite.PostgresPassword, postgres.PostgresPassword)
	suite.Equal(suite.PostgresHost, postgres.PostgresHost)
	suite.Equal(suite.PostgresPort, postgres.PostgresPort)
	suite.Equal(suite.PostgresDatabase, postgres.PostgresDatabase)
	suite.Equal(suite.PostgresSSLMode, postgres.PostgresSSLMode)
	suite.Equal(suite.PostgresURL, postgres.PostgresURL)
	suite.Equal(suite.PostgresTxTimeout, postgres.PostgresTxTimeout)
}

func TestPostgresSuite(t *testing.T) {
	suite.Run(t, new(PostgresSuite))
}
