package errorhandler

import (
	"github.com/cockroachdb/errors"
	"github.com/gogo/status"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"testing"
)

type GRPCHandlerSuite struct {
	suite.Suite
}

func (suite *GRPCHandlerSuite) TestPanicGRPCErrorHandlerNormalStringError() {
	var errContent error
	func() {
		defer PanicGRPCErrorHandler(&errContent, "MockGRPCHandler", "Test error handler")
		panic("got error")
	}()
	suite.Error(errContent)
	if s, ok := status.FromError(errContent); ok {
		suite.Equal("Unknown", s.Code().String())
		suite.Equal("Test error handler: got error", s.Message())
		suite.Equal("rpc error: code = Unknown desc = Test error handler: got error", s.Err().Error())
	}
}

func (suite *GRPCHandlerSuite) TestPanicGRPCErrorHandlerUnknown() {
	var errContent error
	func() {
		defer PanicGRPCErrorHandler(&errContent, "MockGRPCHandler", "Test error handler")
		panic(errors.New("got error"))
	}()
	suite.Error(errContent)
	if s, ok := status.FromError(errContent); ok {
		suite.Equal("Unknown", s.Code().String())
		suite.Equal("Test error handler: got error", s.Message())
		suite.Equal("rpc error: code = Unknown desc = Test error handler: got error", s.Err().Error())
	}
}

func (suite *GRPCHandlerSuite) TestPanicGRPCErrorHandlerStatusError() {
	var errContent error
	func() {
		defer PanicGRPCErrorHandler(&errContent, "MockGRPCHandler", "Test error handler")
		panic(status.Errorf(codes.AlreadyExists, "got error"))
	}()
	suite.Error(errContent)
	if s, ok := status.FromError(errContent); ok {
		suite.Equal("AlreadyExists", s.Code().String())
		suite.Equal("Test error handler: got error", s.Message())
		suite.Equal("rpc error: code = AlreadyExists desc = Test error handler: got error", s.Err().Error())
	}
}

func TestGRPCHandlerSuite(t *testing.T) {
	suite.Run(t, new(GRPCHandlerSuite))
}
