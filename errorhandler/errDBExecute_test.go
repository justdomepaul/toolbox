package errorhandler

import (
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
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

type ErrDBExecuteSuite struct {
	suite.Suite
	obLog *observer.ObservedLogs
}

func (suite *ErrDBExecuteSuite) SetupTest() {
	observedZapCore, observedLogs := observer.New(zap.WarnLevel)
	logger = zap.New(observedZapCore, zap.Fields(zap.String("system", "Mock system")))
	suite.obLog = observedLogs
}

func (suite *ErrDBExecuteSuite) TestNewErrDBExecute() {
	t := suite.T()
	assert.Equal(t, "*errorhandler.ErrDBExecute", reflect.TypeOf(NewErrDBExecute(errors.New("got error"))).String())
}

func (suite *ErrDBExecuteSuite) TestNewErrDBExecuteGetNameMethod() {
	t := suite.T()
	assert.Equal(t, ErrDbExecute, NewErrDBExecute(errors.New("got error")).GetName())
}

func (suite *ErrDBExecuteSuite) TestNewErrDBExecuteGetErrorMethod() {
	t := suite.T()
	assert.Equal(t, errors.New("got error"), NewErrDBExecute(errors.New("got error")).GetError())
}

func (suite *ErrDBExecuteSuite) TestNewErrDBExecuteImplementError() {
	t := suite.T()
	assert.Implements(t, (*error)(nil), NewErrDBExecute(errors.New("got error")))
}

func (suite *ErrDBExecuteSuite) TestNewErrDBExecuteErrorMethod() {
	t := suite.T()
	assert.Equal(t, "[ERROR]: got error\n", NewErrDBExecute(errors.New("got error")).Error())
}

func (suite *ErrDBExecuteSuite) TestNewErrDBExecuteReportMethod() {
	NewErrDBExecute(errors.New("got error")).SetSystem("Mock system").Report("")
	require.Equal(suite.T(), 1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("", firstLog.Message)
	suite.Equal("Mock system", firstLog.Context[0].String)
	suite.Equal("got error", errors.Cause(firstLog.Context[1].Interface.(error)).Error())
}

func (suite *ErrDBExecuteSuite) TestNewErrDBExecuteGinReportMethod() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler("Mock Gin", "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		panic(NewErrDBExecute(errors.New("got error")))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusConflict, result.StatusCode)

	require.Equal(suite.T(), 1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("Mock system", firstLog.Context[0].String)
	suite.Equal("got error", errors.Cause(firstLog.Context[1].Interface.(error)).Error())
}

func (suite *ErrDBExecuteSuite) TestPanicGRPCErrorHandlerNewErrDBExecute() {
	t := suite.T()
	var errContent error
	func() {
		defer PanicGRPCErrorHandler(&errContent, "MockGRPCHandler", "Test error handler")
		panic(NewErrDBExecute(errors.New("database disconnect")))
	}()
	assert.Error(t, errContent)
	if s, ok := status.FromError(errContent); ok {
		assert.Equal(t, "FailedPrecondition", s.Code().String())
		assert.Equal(t, "Test error handler: database disconnect", s.Message())
		assert.Equal(t, "rpc error: code = FailedPrecondition desc = Test error handler: database disconnect", s.Err().Error())
	}
}

func TestErrDBExecuteSuite(t *testing.T) {
	suite.Run(t, new(ErrDBExecuteSuite))
}
