package errorhandler

import (
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/gogo/status"
	"google.golang.org/grpc/codes"
)

type IGRPCErrorReport interface {
	IErrorReport
	GRPCReport(errContent *error, prefixMessage string)
}

func PanicGRPCErrorHandler(errContent *error, system, prefixMessage string) {
	if err := recover(); err != nil {
		switch err.(type) {
		case IGRPCErrorReport:
			err.(IGRPCErrorReport).SetSystem(system).Report(prefixMessage)
			err.(IGRPCErrorReport).GRPCReport(errContent, prefixMessage)
			return
		case error:
			fromStatusError(errContent, err.(error), prefixMessage)
			return
		default:
			*errContent = status.Error(codes.Unknown, errors.Wrap(fmt.Errorf("%v", err), prefixMessage).Error())
		}
	}
}

func fromStatusError(errContent *error, input error, prefixMessage string) {
	if statusErr, ok := status.FromError(input); ok {
		*errContent = status.Error(statusErr.Code(), errors.Wrap(errors.New(statusErr.Message()), prefixMessage).Error())
	} else {
		*errContent = status.Error(codes.Unknown, errors.Wrap(input.(error), prefixMessage).Error())
	}
}
