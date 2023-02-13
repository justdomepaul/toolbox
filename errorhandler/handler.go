package errorhandler

import (
	"fmt"
	"go.uber.org/zap"
	"sync"
)

type IErrorReport interface {
	error
	SetSystem(system string) IErrorReport
	GetName() string
	GetError() error
	Report(prefix string)
}

func PanicErrorHandler(system, prefixMessage string) {
	if err := recover(); err != nil {
		if errReportInstance, okIErrorReport := err.(IErrorReport); okIErrorReport {
			errReportInstance.SetSystem(system).Report("")
			return
		}
		switch err.(type) {
		case error:
			logger.Error(prefixMessage, zap.Error(err.(error)))
		case string:
			logger.Error(fmt.Sprintf("%s: %s", prefixMessage, err))
		default:
			logger.Error(prefixMessage, zap.Any("data", err))
		}
	}
}

func PanicWaitGroupErrorHandler(errContent *error, rmu *sync.RWMutex, wg *sync.WaitGroup) {
	defer wg.Done()
	if err := recover(); err != nil {
		rmu.RLock()
		*errContent = fmt.Errorf("%v", err)
		rmu.RUnlock()
	}
}
