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
	suite.NoError(os.Setenv("PORT", suite.Port))
	suite.NoError(os.Setenv("METRICS_PORT", suite.MetricsPort))
	suite.NoError(os.Setenv("SERVER_TIMEOUT", fmt.Sprint(suite.ServerTimeout)))
	suite.NoError(os.Setenv("PREFIX_MESSAGE", suite.PrefixMessage))
	suite.NoError(os.Setenv("CUSTOMIZED_RENDER", strconv.FormatBool(suite.CustomizedRender)))
	suite.NoError(os.Setenv("ALLOW_ALL_ORIGINS", strconv.FormatBool(suite.AllowAllOrigins)))
	suite.NoError(os.Setenv("ALLOW_ORIGINS", strings.Join(suite.AllowOrigins, ",")))
	suite.NoError(os.Setenv("ALLOWED_PATHS", strings.Join(suite.AllowedPaths, ",")))
	suite.NoError(os.Setenv("ALLOWED_PATHS", strings.Join(suite.AllowedPaths, ",")))
	suite.NoError(os.Setenv("JWT_GUARD", strconv.FormatBool(suite.JWTGuard)))
	suite.NoError(os.Setenv("MAX_MULTIPART_MEMORY_MB", strconv.FormatInt(suite.MaxMultipartMemoryMB, 10)))
}

func (suite *ServerSuite) TestDefaultOption() {
	server := &Server{}
	suite.NoError(LoadFromEnv(server))
	suite.Equal(suite.ReleaseMode, server.ReleaseMode)
	suite.Equal(suite.Port, server.Port)
	suite.Equal(suite.MetricsPort, server.MetricsPort)
	suite.Equal(suite.ServerTimeout, server.ServerTimeout)
	suite.Equal(suite.PrefixMessage, server.PrefixMessage)
	suite.Equal(suite.CustomizedRender, server.CustomizedRender)
	suite.Equal(suite.AllowAllOrigins, server.AllowAllOrigins)
	suite.Equal(suite.AllowOrigins, server.AllowOrigins)
	suite.Equal(suite.AllowedPaths, server.AllowedPaths)
	suite.Equal(suite.JWTGuard, server.JWTGuard)
	suite.Equal(suite.MaxMultipartMemoryMB, server.MaxMultipartMemoryMB)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}
