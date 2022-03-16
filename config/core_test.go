package config

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type CoreSuite struct {
	suite.Suite
	LoggerMode string
	SystemName string
}

func (suite *CoreSuite) SetupSuite() {
	os.Clearenv()
	suite.LoggerMode = "debug"
	suite.SystemName = "system"
	suite.NoError(os.Setenv("LOGGER_MODE", suite.LoggerMode))
	suite.NoError(os.Setenv("SYSTEM_NAME", suite.SystemName))
}

func (suite *CoreSuite) TestDefaultOption() {
	core := &Core{}
	suite.NoError(LoadFromEnv(core))
	suite.Equal(suite.LoggerMode, core.LoggerMode)
	suite.Equal(suite.SystemName, core.SystemName)
}

func TestCoreSuite(t *testing.T) {
	suite.Run(t, new(CoreSuite))
}
