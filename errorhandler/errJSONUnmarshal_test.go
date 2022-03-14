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

type ErrJSONUnmarshalSuite struct {
	suite.Suite
	obLog *observer.ObservedLogs
}

func (suite *ErrJSONUnmarshalSuite) SetupTest() {
	observedZapCore, observedLogs := observer.New(zap.WarnLevel)
	logger = zap.New(observedZapCore, zap.Fields(zap.String("system", "Mock system")))
	suite.obLog = observedLogs
}

func (suite *ErrJSONUnmarshalSuite) TestNewErrJSONUnmarshal() {
	suite.Equal("*errorhandler.ErrJSONUnmarshal", reflect.TypeOf(NewErrJSONUnmarshal(errors.New("got error"))).String())
}

func (suite *ErrJSONUnmarshalSuite) TestNewErrJSONUnmarshalGetNameMethod() {
	suite.Equal(ErrJsonUnmarshal, NewErrJSONUnmarshal(errors.New("got error")).GetName())
}

func (suite *ErrJSONUnmarshalSuite) TestNewErrJSONUnmarshalGetErrorMethod() {
	suite.Equal(errors.New("got error"), NewErrJSONUnmarshal(errors.New("got error")).GetError())
}

func (suite *ErrJSONUnmarshalSuite) TestNewErrJSONUnmarshalImplementError() {
	suite.Implements((*error)(nil), NewErrJSONUnmarshal(errors.New("got error")))
}

func (suite *ErrJSONUnmarshalSuite) TestNewErrJSONUnmarshalErrorMethod() {
	suite.Equal("[ERROR]: got error\n", NewErrJSONUnmarshal(errors.New("got error")).Error())
}

func (suite *ErrJSONUnmarshalSuite) TestNewErrJSONUnmarshalReportMethod() {
	NewErrJSONUnmarshal(errors.New("got error")).SetSystem("Mock system").Report("Error JSON")
	require.Equal(suite.T(), 1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("Error JSON", firstLog.Message)
	suite.Equal("Mock system", firstLog.Context[0].String)
	suite.Equal("got error", errors.Cause(firstLog.Context[1].Interface.(error)).Error())
}

func (suite *ErrJSONUnmarshalSuite) TestNewErrJSONUnmarshalGinReportMethod() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.New()
	route.Use(gin.Logger(), GinPanicErrorHandler("Mock Gin", "error Gin mock"))
	route.GET("/", func(c *gin.Context) {
		panic(NewErrJSONUnmarshal(errors.New("got error")))
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	route.ServeHTTP(w, req)
	result := w.Result()
	defer result.Body.Close()
	suite.Equal(http.StatusBadRequest, result.StatusCode)

	require.Equal(suite.T(), 1, suite.obLog.Len())
	firstLog := suite.obLog.All()[0]
	suite.Equal("error Gin mock", firstLog.Message)
	suite.Equal("Mock system", firstLog.Context[0].String)
	suite.Equal("got error", errors.Cause(firstLog.Context[1].Interface.(error)).Error())
}

func (suite *ErrJSONUnmarshalSuite) TestPanicGRPCErrorHandlerNewErrJSONMarshal() {
	var errContent error
	func() {
		defer PanicGRPCErrorHandler(&errContent, "MockGRPCHandler", "Test error handler")
		panic(NewErrJSONUnmarshal(errors.New("json unmarshal disconnect")))
	}()
	suite.Error(errContent)
	if s, ok := status.FromError(errContent); ok {
		suite.Equal("InvalidArgument", s.Code().String())
		suite.Equal("Test error handler: json unmarshal disconnect", s.Message())
		suite.Equal("rpc error: code = InvalidArgument desc = Test error handler: json unmarshal disconnect", s.Err().Error())
	}
}

func TestErrJSONUnmarshalSuite(t *testing.T) {
	suite.Run(t, new(ErrJSONUnmarshalSuite))
}
