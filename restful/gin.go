package restful

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/justdomepaul/toolbox/config"
	"github.com/justdomepaul/toolbox/errorhandler"
	"html/template"
	"net/http"
)

func NewGin(
	option config.Set,
	render *Render,
	guarder *JWTGuarder,
) (*gin.Engine, error) {
	if option.Server.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	srv := gin.New()

	if option.Server.CustomizedRender {
		tmpl, err := template.New("tmpl").Parse(ErrorPageTmpl)
		if err != nil {
			return nil, err
		}
		if err := render.Add(ErrPageKey, tmpl); err != nil {
			return nil, err
		}
		srv.HTMLRender = render
	}

	srv.MaxMultipartMemory = option.Server.MaxMultipartMemoryMB << 20

	cf := cors.DefaultConfig()
	cf.AllowMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions}
	cf.AllowHeaders = []string{
		"Origin",
		"Upgrade",
		"Content-Length",
		"Content-Type",
		"Authorization",
		"Connection",
		"Accept-Encoding",
		"Accept-Language",
		"Host",
		"X-Google-*",
		"X-AppEngine-*",
		"X-CloudScheduler",
		"X-CloudScheduler-JobName",
		"X-CloudScheduler-ScheduleTime",
		"Sec-WebSocket-Key",
		"Sec-WebSocket-Version",
		"Sec-WebSocket-Protocol",
	}
	if option.Server.AllowAllOrigins {
		cf.AllowAllOrigins = true
	} else {
		cf.AllowOrigins = option.Server.AllowOrigins
	}
	fns := []gin.HandlerFunc{
		cors.New(cf),
		gin.Logger(),
		errorhandler.GinPanicErrorHandler(option.Core.SystemName, option.Server.PrefixMessage),
	}
	if option.Server.JWTGuard {
		fns = append(fns, guarder.JWTGuarder(option.Server.AllowedPaths...))
	}
	srv.Use(fns...)

	return srv, nil
}
