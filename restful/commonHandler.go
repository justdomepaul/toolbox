package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	Error404Set    = Error404()
	QuickReplySet  = QuickReply()
	NewPromHTTPSet = NewPromHTTP()
)

// Error404 method
func Error404() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusNotFound, ErrPageKey, ErrPage{Title: "Page Not Found", Code: http.StatusNotFound})
	}
}

// QuickReply method
func QuickReply() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
		c.Writer.Flush()
	}
}

func NewPromHTTP() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}

type CommonHandler struct {
	Error404   gin.HandlerFunc
	QuickReply gin.HandlerFunc
	PromHTTP   gin.HandlerFunc
}
