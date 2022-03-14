package errorhandler

import (
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"google.golang.org/grpc/status"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type ErrAuthenticateSuite struct {
	suite.Suite
	obLog *observer.ObservedLogs
}

func (suite *ErrAuthenticateSuite) SetupTest() {
	observedZapCore, observedLogs := observer.New(zap.WarnLevel)
	logger = zap.New(observedZapCore, zap.Fields(zap.String("system", "Mock system")))
	suite.obLog = observedLogs
}

func (suite *ErrAuthenticateSuite) TestNewErrAuthenticate() {
	suite.Equal("*errorhandler.ErrAuthenticate", reflect.TypeOf(NewErrAuthenticate(errors.New("got error"))).String())
}

func (suite *ErrAuthenticateSuite) TestNewErrAuthenticateGetNameMethod() {
	suite.Equal(ErrProcessAuthenticate, NewErrAuthenticate(errors.New("got error")).GetName())
}

func (suite *ErrAuthenticateSuite) TestNewErrAuthenticateGetErrorMethod() {
	suite.Equal(errors.New("got error"), NewErrAuthenticate(errors.New("got error")).GetError())
}

func (suite *ErrAuthenticateSuite) TestNewErrAuthenticateImplementError() {
	suite.Implements((*error)(nil), NewErrAuthenticate(errors.New("got error")))
}

func (suite *ErrAuthenticateSuite) TestNewErrAuthenticateErrorMethod() {
	suite.Equal("[ERROR]: got error\n", NewErrAuthenticate(errors.New("got error")).Error())
}

func (suite *ErrAuthenticateSuite) TestNewErrAuthenticateReportMethod() {
	NewErrAuthenticate(errors.New("got error")).SetSystem("Mock system").Report("")
	require.Equal(suite.T(), 1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("", firstLog.Message)
	suite.Equal("Mock system", firstLog.Context[0].String)
	suite.Equal("got error", errors.Cause(firstLog.Context[1].Interface.(error)).Error())
}

func (suite *ErrAuthenticateSuite) TestNewErrAuthenticateGinReportMethod() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler("Mock Gin", "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		panic(NewErrAuthenticate(errors.New("got error")))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusUnauthorized, result.StatusCode)

	require.Equal(suite.T(), 1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("Mock system", firstLog.Context[0].String)
	suite.Equal("got error", errors.Cause(firstLog.Context[1].Interface.(error)).Error())
}

func (suite *ErrAuthenticateSuite) TestPanicGRPCErrorHandlerNewErrAuthenticate() {
	var errContent error
	func() {
		defer PanicGRPCErrorHandler(&errContent, "MockGRPCHandler", "Test error handler")
		panic(NewErrAuthenticate(errors.New("database disconnect")))
	}()
	suite.Error(errContent)
	if s, ok := status.FromError(errContent); ok {
		suite.Equal("Unauthenticated", s.Code().String())
		suite.Equal("Test error handler: database disconnect", s.Message())
		suite.Equal("rpc error: code = Unauthenticated desc = Test error handler: database disconnect", s.Err().Error())
	}
}

func TestErrAuthenticateSuite(t *testing.T) {
	suite.Run(t, new(ErrAuthenticateSuite))
}
