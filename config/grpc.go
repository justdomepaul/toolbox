package config

import (
	"time"
)

// GRPC type
type GRPC struct {
	Port                         string        `split_words:"true" default:"38080"`
	NoTLS                        bool          `split_words:"true" default:"true"`  // use WithInsecure option
	SkipTLS                      bool          `split_words:"true" default:"false"` // use WithTransportCredentials option and include InsecureSkipVerify=true
	TLS                          bool          `split_words:"true" default:"false"` // use WithTransportCredentials option and include PemCert
	TLSNil                       bool          `split_words:"true" default:"false"` // use WithTransportCredentials option and include Use RootCA
	TLSPemCert                   string        `split_words:"true" default:"false"` // TLS PemCert data
	TLSPemCertBase64             string        `split_words:"true" default:"false"` // TLS PemCert base64encode data
	ALTS                         bool          `split_words:"true" default:"false"` // client use grpc.WithTransportCredentials(alts.NewClientCreds(alts.DefaultClientOptions())); server use grpc.Creds(alts.NewServerCreds(alts.DefaultServerOptions()))
	KeepAliveTime                time.Duration `split_words:"true" default:"1s"`    // second
	KeepAliveTimeout             time.Duration `split_words:"true" default:"20s"`   //second
	KeepAlivePermitWithoutStream bool          `split_words:"true" default:"true"`
	AllowedList                  []string      `split_words:"true" default:"/auth.Auth/Ping,/auth.Auth/Authorization"`
}
