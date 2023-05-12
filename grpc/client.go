package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"github.com/cockroachdb/errors"
	"github.com/justdomepaul/toolbox/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/alts"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

var (
	// Dial variable
	Dial = grpc.Dial
)

// IRegister interface
type IRegister interface {
	Handler(*grpc.Server)
}

// IService interface
type IService interface {
	GetSession() (IClientConn, error)
	Close() error
}

// IClientConn interface
type IClientConn interface {
	WaitForStateChange(ctx context.Context, sourceState connectivity.State) bool
	GetState() connectivity.State
	Target() string
	GetMethodConfig(method string) grpc.MethodConfig
	ResetConnectBackoff()
	Close() error
	NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error)
	Invoke(ctx context.Context, method string, args, reply interface{}, opts ...grpc.CallOption) error
}

// CreateClient method
func CreateClient(domain string, option config.GRPC) (IClientConn, error) {
	options := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                option.KeepAliveTime,
			Timeout:             option.KeepAliveTimeout,
			PermitWithoutStream: option.KeepAlivePermitWithoutStream,
		}),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	}
	if option.NoTLS {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	if option.SkipTLS {
		options = append(options,
			grpc.WithTransportCredentials(
				credentials.NewTLS(&tls.Config{
					InsecureSkipVerify: true,
				}),
			),
		)
	}
	if option.TLSNil {
		options = append(options,
			grpc.WithTransportCredentials(credentials.NewTLS(nil)),
		)
	}
	if option.TLS {
		if option.TLSPemCert == "" && option.TLSPemCertBase64 == "" {
			return nil, errors.New("gRPC TLS Pem Cert Data Not Found")
		}
		if option.TLSPemCert != "" && option.TLSPemCertBase64 == "" {
			cp := x509.NewCertPool()
			if !cp.AppendCertsFromPEM([]byte(option.TLSPemCert)) {
				return nil, errors.New("gRPC TLS Pem Cert format error")
			}
			options = append(options,
				grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(cp, "")),
			)
		}
		if option.TLSPemCert == "" && option.TLSPemCertBase64 != "" {
			tlsPemCert, err := base64.StdEncoding.DecodeString(option.TLSPemCertBase64)
			if err != nil {
				return nil, err
			}
			cp := x509.NewCertPool()
			if !cp.AppendCertsFromPEM(tlsPemCert) {
				return nil, errors.New("gRPC TLS Pem Cert format error")
			}
			options = append(options,
				grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(cp, "")),
			)
		}
	}
	if option.ALTS {
		options = append(options, grpc.WithTransportCredentials(alts.NewClientCreds(alts.DefaultClientOptions())))
	}

	return Dial(
		domain,
		options...,
	)
}
