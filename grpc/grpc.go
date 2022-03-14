package grpc

import (
	"context"
	"crypto/tls"
	"github.com/justdomepaul/toolbox/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

var (
	// Dial variable
	Dial = grpc.Dial
)

// Register interface
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
		options = append(options, grpc.WithInsecure())
	}
	if option.SkipTLS {
		cg := &tls.Config{
			InsecureSkipVerify: true,
		}
		cert := credentials.NewTLS(cg)
		options = append(options, grpc.WithTransportCredentials(cert))
	}

	return Dial(
		domain,
		options...,
	)
}
