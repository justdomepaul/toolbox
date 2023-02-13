package errorhandler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type IGinErrorReport interface {
	IErrorReport
	GinReport(c *gin.Context)
}

func GinPanicErrorHandler(system, prefixMessage string) func(c *gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				switch err.(type) {
				case IGinErrorReport:
					err.(IGinErrorReport).SetSystem(system).Report(prefixMessage)
					err.(IGinErrorReport).GinReport(c)
					return
				case error:
					c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("%s: %w", prefixMessage, err.(error)))
					return
				default:
					c.AbortWithError(http.StatusInternalServerError, fmt.Errorf("%s: %s", prefixMessage, err))
				}
			}
		}()
		c.Next()
	}
}
