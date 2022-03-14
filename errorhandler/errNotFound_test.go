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

type ErrNotFoundSuite struct {
	suite.Suite
	obLog *observer.ObservedLogs
}

func (suite *ErrNotFoundSuite) SetupTest() {
	observedZapCore, observedLogs := observer.New(zap.WarnLevel)
	logger = zap.New(observedZapCore, zap.Fields(zap.String("system", "Mock system")))
	suite.obLog = observedLogs
}

func (suite *ErrNotFoundSuite) TestNewErrNotFound() {
	suite.Equal("*errorhandler.ErrNotFound", reflect.TypeOf(NewErrNotFound(errors.New("got error"))).String())
}

func (suite *ErrNotFoundSuite) TestNewErrNotFoundGetNameMethod() {
	suite.Equal(ErrDataNotFound, NewErrNotFound(errors.New("got error")).GetName())
}

func (suite *ErrNotFoundSuite) TestNewErrNotFoundGetErrorMethod() {
	suite.Equal(errors.New("got error"), NewErrNotFound(errors.New("got error")).GetError())
}

func (suite *ErrNotFoundSuite) TestNewErrNotFoundImplementError() {
	suite.Implements((*error)(nil), NewErrNotFound(errors.New("got error")))
}

func (suite *ErrNotFoundSuite) TestNewErrNotFoundErrorMethod() {
	suite.Equal("[ERROR]: got error\n", NewErrNotFound(errors.New("got error")).Error())
}

func (suite *ErrNotFoundSuite) TestNewErrNotFoundReportMethod() {
	NewErrNotFound(errors.New("got error")).SetSystem("Mock system").Report("")
	require.Equal(suite.T(), 1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("", firstLog.Message)
	suite.Equal("Mock system", firstLog.Context[0].String)
	suite.Equal("got error", errors.Cause(firstLog.Context[1].Interface.(error)).Error())
}

func (suite *ErrNotFoundSuite) TestNewErrNotFoundGinReportMethod() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler("Mock Gin", "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		panic(NewErrNotFound(errors.New("got error")))
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

func (suite *ErrNotFoundSuite) TestPanicGRPCErrorHandlerNewErrNotFound() {
	var errContent error
	func() {
		defer PanicGRPCErrorHandler(&errContent, "MockGRPCHandler", "Test error handler")
		panic(NewErrNotFound(errors.New("database disconnect")))
	}()
	suite.Error(errContent)
	if s, ok := status.FromError(errContent); ok {
		suite.Equal("NotFound", s.Code().String())
		suite.Equal("Test error handler: database disconnect", s.Message())
		suite.Equal("rpc error: code = NotFound desc = Test error handler: database disconnect", s.Err().Error())
	}
}

func TestErrNotFoundSuite(t *testing.T) {
	suite.Run(t, new(ErrNotFoundSuite))
}
