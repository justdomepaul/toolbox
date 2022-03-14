package config

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

type CockroachSuite struct {
	suite.Suite
	CockroachUsername  string
	CockroachPassword  string
	CockroachHost      string
	CockroachPort      string
	CockroachDatabase  string
	CockroachSSLMode   string
	CockroachURL       string
	CockroachTxTimeout time.Duration
}

func (suite *CockroachSuite) SetupSuite() {
	os.Clearenv()
	suite.CockroachUsername = "testCockroachUsername"
	suite.CockroachPassword = "testCockroachPassword"
	suite.CockroachHost = "testCockroachHost"
	suite.CockroachPort = "testCockroachPort"
	suite.CockroachDatabase = "testCockroachDatabase"
	suite.CockroachSSLMode = "testCockroachSSLMode"
	suite.CockroachURL = "testCockroachURL"
	suite.CockroachTxTimeout = 20 * time.Second
	suite.NoError(os.Setenv("COCKROACH_USERNAME", suite.CockroachUsername))
	suite.NoError(os.Setenv("COCKROACH_PASSWORD", suite.CockroachPassword))
	suite.NoError(os.Setenv("COCKROACH_HOST", suite.CockroachHost))
	suite.NoError(os.Setenv("COCKROACH_PORT", suite.CockroachPort))
	suite.NoError(os.Setenv("COCKROACH_DATABASE", suite.CockroachDatabase))
	suite.NoError(os.Setenv("COCKROACH_SSL_MODE", suite.CockroachSSLMode))
	suite.NoError(os.Setenv("COCKROACH_URL", suite.CockroachURL))
	suite.NoError(os.Setenv("COCKROACH_TX_TIMEOUT", fmt.Sprint(suite.CockroachTxTimeout)))
}

func (suite *CockroachSuite) TestDefaultOption() {
	cockroach := &Cockroach{}
	suite.NoError(LoadFromEnv(cockroach))
	suite.Equal(suite.CockroachUsername, cockroach.CockroachUsername)
	suite.Equal(suite.CockroachPassword, cockroach.CockroachPassword)
	suite.Equal(suite.CockroachHost, cockroach.CockroachHost)
	suite.Equal(suite.CockroachPort, cockroach.CockroachPort)
	suite.Equal(suite.CockroachDatabase, cockroach.CockroachDatabase)
	suite.Equal(suite.CockroachSSLMode, cockroach.CockroachSSLMode)
	suite.Equal(suite.CockroachURL, cockroach.CockroachURL)
	suite.Equal(suite.CockroachTxTimeout, cockroach.CockroachTxTimeout)
}

func TestCockroachSuite(t *testing.T) {
	suite.Run(t, new(CockroachSuite))
}
