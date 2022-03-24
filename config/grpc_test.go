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

type GRPCSuite struct {
	suite.Suite
	Port                         string
	NoTLS                        bool
	SkipTLS                      bool
	ALTS                         bool
	KeepAliveTime                time.Duration
	KeepAliveTimeout             time.Duration
	KeepAlivePermitWithoutStream bool
	AllowedList                  []string
}

func (suite *GRPCSuite) SetupSuite() {
	os.Clearenv()
	suite.Port = "Port"
	suite.NoTLS = true
	suite.SkipTLS = false
	suite.ALTS = false
	suite.KeepAliveTime = 200 * time.Second
	suite.KeepAliveTimeout = 200 * time.Second
	suite.KeepAlivePermitWithoutStream = true
	suite.AllowedList = []string{"a", "b", "c"}

	suite.NoError(os.Setenv("PORT", suite.Port))
	suite.NoError(os.Setenv("NO_TLS", strconv.FormatBool(suite.NoTLS)))
	suite.NoError(os.Setenv("SKIP_TLS", strconv.FormatBool(suite.SkipTLS)))
	suite.NoError(os.Setenv("ALTS", strconv.FormatBool(suite.ALTS)))
	suite.NoError(os.Setenv("KEEP_ALIVE_TIME", fmt.Sprint(suite.KeepAliveTime)))
	suite.NoError(os.Setenv("KEEP_ALIVE_TIMEOUT", fmt.Sprint(suite.KeepAliveTimeout)))
	suite.NoError(os.Setenv("KEEP_ALIVE_PERMIT_WITHOUT_STREAM", strconv.FormatBool(suite.KeepAlivePermitWithoutStream)))
	suite.NoError(os.Setenv("ALLOWED_LIST", strings.Join(suite.AllowedList, ",")))
}

func (suite *GRPCSuite) TestDefaultOption() {
	grpc := &GRPC{}
	suite.NoError(LoadFromEnv(grpc))
	suite.Equal(suite.Port, grpc.Port)
	suite.Equal(suite.NoTLS, grpc.NoTLS)
	suite.Equal(suite.SkipTLS, grpc.SkipTLS)
	suite.Equal(suite.ALTS, grpc.ALTS)
	suite.Equal(suite.KeepAliveTime, grpc.KeepAliveTime)
	suite.Equal(suite.KeepAliveTimeout, grpc.KeepAliveTimeout)
	suite.Equal(suite.KeepAlivePermitWithoutStream, grpc.KeepAlivePermitWithoutStream)
	suite.Equal(suite.AllowedList, grpc.AllowedList)
}

func TestGRPCSuite(t *testing.T) {
	suite.Run(t, new(GRPCSuite))
}
