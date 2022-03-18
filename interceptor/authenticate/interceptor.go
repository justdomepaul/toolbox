package authenticate

import (
	"context"
	"github.com/cockroachdb/errors"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/justdomepaul/toolbox/definition"
	"github.com/justdomepaul/toolbox/errorhandler"
	"github.com/justdomepaul/toolbox/services"
	"github.com/justdomepaul/toolbox/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor(auth services.Authenticate) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		id, err := authenticate(ctx, auth, info.FullMethod)
		if err != nil {
			return nil, err
		}
		return handler(utils.SetClientID(ctx, definition.AuthorizationClientKey, id), req)
	}
}

func StreamServerInterceptor(auth services.Authenticate) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		id, err := authenticate(ss.Context(), auth, info.FullMethod)
		if err != nil {
			return err
		}
		wrapped := grpc_middleware.WrapServerStream(ss)
		wrapped.WrappedContext = utils.SetClientID(ss.Context(), definition.AuthorizationClientKey, id)
		return handler(srv, wrapped)
	}
}

func authenticate(ctx context.Context, auth services.Authenticate, fullMethod string) ([]byte, error) {
	id, err := auth.Authenticate(nil, func() (string, error) {
		return utils.GetAccessToken(ctx)
	}, fullMethod)
	if _, exist := status.FromError(err); err != nil && exist {
		return id, err
	}
	if errors.Is(err, errorhandler.ErrInWhitelist) {
		return id, nil
	}
	if errors.Is(err, errorhandler.ErrUnauthenticated) {
		return id, status.Errorf(codes.Unauthenticated, "%v", err)
	}
	if errors.Is(err, errorhandler.ErrNoPermission) {
		return id, status.Errorf(codes.PermissionDenied, "%v", err)
	}
	if errors.Is(err, errorhandler.ErrInvalidArguments) {
		return id, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	if err != nil {
		return id, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	return id, nil
}
