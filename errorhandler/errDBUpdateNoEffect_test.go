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

type ErrDBUpdateNoEffectSuite struct {
	suite.Suite
	obLog *observer.ObservedLogs
}

func (suite *ErrDBUpdateNoEffectSuite) SetupTest() {
	observedZapCore, observedLogs := observer.New(zap.WarnLevel)
	logger = zap.New(observedZapCore, zap.Fields(zap.String("system", "Mock system")))
	suite.obLog = observedLogs
}

func (suite *ErrDBUpdateNoEffectSuite) TestNewErrDBUpdateNoEffect() {
	suite.Equal("*errorhandler.ErrDBUpdateNoEffect", reflect.TypeOf(NewErrDBUpdateNoEffect(errors.New("got error"))).String())
}

func (suite *ErrDBUpdateNoEffectSuite) TestNewErrDBUpdateNoEffectGetNameMethod() {
	suite.Equal(ErrDbUpdateNoEffect, NewErrDBUpdateNoEffect(errors.New("got error")).GetName())
}

func (suite *ErrDBUpdateNoEffectSuite) TestNewErrDBUpdateNoEffectGetErrorMethod() {
	suite.Equal(errors.New("got error"), NewErrDBUpdateNoEffect(errors.New("got error")).GetError())
}

func (suite *ErrDBUpdateNoEffectSuite) TestNewErrDBUpdateNoEffectImplementError() {
	suite.Implements((*error)(nil), NewErrDBUpdateNoEffect(errors.New("got error")))
}

func (suite *ErrDBUpdateNoEffectSuite) TestNewErrDBUpdateNoEffectErrorMethod() {
	suite.Equal("[ERROR]: got error\n", NewErrDBUpdateNoEffect(errors.New("got error")).Error())
}

func (suite *ErrDBUpdateNoEffectSuite) TestNewErrDBUpdateNoEffectReportMethod() {
	NewErrDBUpdateNoEffect(errors.New("got error")).SetSystem("Mock system").Report("")
	require.Equal(suite.T(), 1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("", firstLog.Message)
	suite.Equal("Mock system", firstLog.Context[0].String)
	suite.Equal("got error", errors.Cause(firstLog.Context[1].Interface.(error)).Error())
}

func (suite *ErrDBUpdateNoEffectSuite) TestNewErrDBUpdateNoEffectGinReportMethod() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler("Mock Gin", "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		panic(NewErrDBUpdateNoEffect(errors.New("got error")))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusNotFound, result.StatusCode)

	require.Equal(suite.T(), 1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("Mock system", firstLog.Context[0].String)
	suite.Equal("got error", errors.Cause(firstLog.Context[1].Interface.(error)).Error())
}

func (suite *ErrDBUpdateNoEffectSuite) TestPanicGRPCErrorHandlerNewErrDBUpdateNoEffect() {
	var errContent error
	func() {
		defer PanicGRPCErrorHandler(&errContent, "MockGRPCHandler", "Test error handler")
		panic(NewErrDBUpdateNoEffect(errors.New("database disconnect")))
	}()
	suite.Error(errContent)
	if s, ok := status.FromError(errContent); ok {
		suite.Equal("NotFound", s.Code().String())
		suite.Equal("Test error handler: database disconnect", s.Message())
		suite.Equal("rpc error: code = NotFound desc = Test error handler: database disconnect", s.Err().Error())
	}
}

func TestErrDBUpdateNoEffectSuite(t *testing.T) {
	suite.Run(t, new(ErrDBUpdateNoEffectSuite))
}
