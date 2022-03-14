package config

import (
	"github.com/stretchr/testify/suite"
	"os"
	"strconv"
	"testing"
)

type CoreSuite struct {
	suite.Suite
	LoggerMode  string
	ReleaseMode bool
	SystemName  string
}

func (suite *CoreSuite) SetupSuite() {
	os.Clearenv()
	suite.LoggerMode = "debug"
	suite.ReleaseMode = true
	suite.SystemName = "system"
	suite.NoError(os.Setenv("LOGGER_MODE", suite.LoggerMode))
	suite.NoError(os.Setenv("RELEASE_MODE", strconv.FormatBool(suite.ReleaseMode)))
	suite.NoError(os.Setenv("SYSTEM_NAME", suite.SystemName))
}

func (suite *CoreSuite) TestDefaultOption() {
	core := &Core{}
	suite.NoError(LoadFromEnv(core))
	suite.Equal(suite.LoggerMode, core.LoggerMode)
	suite.Equal(suite.ReleaseMode, core.ReleaseMode)
	suite.Equal(suite.SystemName, core.SystemName)
}

func TestCoreSuite(t *testing.T) {
	suite.Run(t, new(CoreSuite))
}
