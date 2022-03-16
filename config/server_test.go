package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

type ServerSuite struct {
	suite.Suite
	ReleaseMode          bool
	Port                 string
	MetricsPort          string
	ServerTimeout        time.Duration
	PrefixMessage        string
	CustomizedRender     bool
	AllowAllOrigins      bool
	AllowOrigins         []string
	AllowedPaths         []string
	JWTGuard             bool
	MaxMultipartMemoryMB int64
}

func (suite *ServerSuite) SetupSuite() {
	t := suite.T()
	os.Clearenv()
	suite.ReleaseMode = true
	suite.Port = "Port"
	suite.MetricsPort = "Port"
	suite.ServerTimeout = 5 * time.Second
	suite.PrefixMessage = "testPrefixMessage"
	suite.CustomizedRender = true
	suite.AllowAllOrigins = true
	suite.AllowOrigins = []string{"testAllowOrigins", "testAllowOrigins2"}
	suite.AllowedPaths = []string{"/user/v1/login", "/user/v1/logout", "/user/v1/refresh_token"}
	suite.JWTGuard = false
	suite.MaxMultipartMemoryMB = 16

	suite.NoError(os.Setenv("RELEASE_MODE", strconv.FormatBool(suite.ReleaseMode)))
	assert.NoError(t, os.Setenv("PORT", suite.Port))
	assert.NoError(t, os.Setenv("METRICS_PORT", suite.MetricsPort))
	assert.NoError(t, os.Setenv("SERVER_TIMEOUT", fmt.Sprint(suite.ServerTimeout)))
	assert.NoError(t, os.Setenv("PREFIX_MESSAGE", suite.PrefixMessage))
	assert.NoError(t, os.Setenv("CUSTOMIZED_RENDER", strconv.FormatBool(suite.CustomizedRender)))
	assert.NoError(t, os.Setenv("ALLOW_ALL_ORIGINS", strconv.FormatBool(suite.AllowAllOrigins)))
	assert.NoError(t, os.Setenv("ALLOW_ORIGINS", strings.Join(suite.AllowOrigins, ",")))
	assert.NoError(t, os.Setenv("ALLOWED_PATHS", strings.Join(suite.AllowedPaths, ",")))
	assert.NoError(t, os.Setenv("ALLOWED_PATHS", strings.Join(suite.AllowedPaths, ",")))
	assert.NoError(t, os.Setenv("JWT_GUARD", strconv.FormatBool(suite.JWTGuard)))
	assert.NoError(t, os.Setenv("MAX_MULTIPART_MEMORY_MB", strconv.FormatInt(suite.MaxMultipartMemoryMB, 10)))
}

func (suite *ServerSuite) TestDefaultOption() {
	t := suite.T()
	server := &Server{}
	suite.NoError(LoadFromEnv(server))
	suite.Equal(suite.ReleaseMode, server.ReleaseMode)
	assert.Equal(t, suite.Port, server.Port)
	assert.Equal(t, suite.MetricsPort, server.MetricsPort)
	assert.Equal(t, suite.ServerTimeout, server.ServerTimeout)
	assert.Equal(t, suite.PrefixMessage, server.PrefixMessage)
	assert.Equal(t, suite.CustomizedRender, server.CustomizedRender)
	assert.Equal(t, suite.AllowAllOrigins, server.AllowAllOrigins)
	assert.Equal(t, suite.AllowOrigins, server.AllowOrigins)
	assert.Equal(t, suite.AllowedPaths, server.AllowedPaths)
	assert.Equal(t, suite.JWTGuard, server.JWTGuard)
	assert.Equal(t, suite.MaxMultipartMemoryMB, server.MaxMultipartMemoryMB)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
