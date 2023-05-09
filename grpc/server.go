package grpc

import (
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/justdomepaul/toolbox/config"
	"github.com/justdomepaul/toolbox/interceptor/authenticate"
	"github.com/justdomepaul/toolbox/services"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/alts"
	"google.golang.org/grpc/keepalive"
	"time"
)

func CreateServer(logger *zap.Logger, grpcOption config.GRPC, authenticateService services.IAuthenticate) *grpc.Server {
	opts := []grpc_zap.Option{
		grpc_zap.WithDurationField(func(duration time.Duration) zapcore.Field {
			return zap.Int64("grpc.time_ns", duration.Nanoseconds())
		}),
	}
	options := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_zap.UnaryServerInterceptor(logger, opts...),
			grpc_prometheus.UnaryServerInterceptor,
			authenticate.UnaryServerInterceptor(authenticateService),
		),
		grpc.ChainStreamInterceptor(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(logger, opts...),
			grpc_prometheus.StreamServerInterceptor,
			authenticate.StreamServerInterceptor(authenticateService),
		),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Timeout: grpcOption.KeepAliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcOption.KeepAliveTime,
			PermitWithoutStream: grpcOption.KeepAlivePermitWithoutStream,
		}),
	}
	if grpcOption.ALTS {
		options = append(options, grpc.Creds(alts.NewServerCreds(alts.DefaultServerOptions())))
	}

	return grpc.NewServer(
		options...,
	)
}
