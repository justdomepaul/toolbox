package config

import (
	"time"
)

// GRPC type
type GRPC struct {
	Port                         string        `split_words:"true" default:"38080"`
	NoTLS                        bool          `split_words:"true" default:"true"`  // use WithInsecure option
	SkipTLS                      bool          `split_words:"true" default:"false"` // use WithTransportCredentials option and include InsecureSkipVerify=true
	KeepAliveTime                time.Duration `split_words:"true" default:"1s"`    // second
	KeepAliveTimeout             time.Duration `split_words:"true" default:"20s"`   //second
	KeepAlivePermitWithoutStream bool          `split_words:"true" default:"true"`
	AllowedList                  []string      `split_words:"true" default:"/auth.Auth/Ping,/auth.Auth/Authorization"`
}
