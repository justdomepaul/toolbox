package errorhandler

import (
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

type ErrGRPCExecute struct {
	system string
	err    error
}

func (e *ErrGRPCExecute) SetSystem(system string) IErrorReport {
	if e.system == "" {
		e.system = system
	}
	return e
}

// GetName method
func (e ErrGRPCExecute) GetName() string {
	return ErrGrpcExecute
}

func (e ErrGRPCExecute) GetError() error {
	return e.err
}

func (e ErrGRPCExecute) Error() string {
	return fmt.Sprintln("[ERROR]:", e.err.Error())
}

func (e ErrGRPCExecute) Report(prefix string) {
	logger.Warn(prefix, zap.Error(e.GetError()))
}

func (e ErrGRPCExecute) GinReport(c *gin.Context) {
	c.AbortWithError(http.StatusUnprocessableEntity, e.err)
}

func (e ErrGRPCExecute) GRPCReport(errContent *error, prefixMessage string) {
	*errContent = status.Error(codes.FailedPrecondition, errors.Wrap(e.err, prefixMessage).Error())
}

func NewErrGRPCExecute(err error) *ErrGRPCExecute {
	return &ErrGRPCExecute{
		err: err,
	}
}
