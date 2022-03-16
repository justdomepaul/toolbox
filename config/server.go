package config

import (
	"time"
)

// Server type
type Server struct {
	ReleaseMode          bool          `split_words:"true" default:"true"`
	Port                 string        `split_words:"true" default:"38080"`
	MetricsPort          string        `split_words:"true" default:"38090"`
	ServerTimeout        time.Duration `split_words:"true" default:"5s"`
	PrefixMessage        string        `split_words:"true" default:"error gin server"`
	CustomizedRender     bool          `split_words:"true" default:"false"`
	AllowAllOrigins      bool          `split_words:"true" default:"false"`
	AllowOrigins         []string      `split_words:"true" default:"http://localhost,https://localhost"`
	AllowedPaths         []string      `split_words:"true" default:"/favicon.ico,/ping,/metrics,/api/auth/v1/authorization,/narrow_cast_schedule"`
	JWTGuard             bool          `split_words:"true" default:"true"`
	MaxMultipartMemoryMB int64         `split_words:"true" default:"8"`
}
