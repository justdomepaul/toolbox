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

func UnaryServerInterceptor(auth services.IAuthenticate) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		result, err := authenticate(ctx, auth, info.FullMethod)
		if err != nil {
			return nil, err
		}
		return handler(
			utils.SetClaim(
				utils.SetID(
					ctx,
					definition.AuthorizationID,
					result.GetID(),
				),
				definition.AuthorizationClaim,
				result.GetClaim(),
			),
			req,
		)
	}
}

func StreamServerInterceptor(auth services.IAuthenticate) func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		result, err := authenticate(ss.Context(), auth, info.FullMethod)
		if err != nil {
			return err
		}
		wrapped := grpc_middleware.WrapServerStream(ss)
		wrapped.WrappedContext = utils.SetClaim(
			utils.SetID(
				ss.Context(),
				definition.AuthorizationID,
				result.GetID(),
			),
			definition.AuthorizationClaim,
			result.GetClaim(),
		)
		return handler(srv, wrapped)
	}
}

func authenticate(ctx context.Context, auth services.IAuthenticate, fullMethod string) (services.IAuthorization, error) {
	result, err := auth.Authenticate(nil, func() (string, error) {
		return utils.GetAccessToken(ctx)
	}, fullMethod)
	if _, exist := status.FromError(err); err != nil && exist {
		return result, err
	}
	if errors.Is(err, errorhandler.ErrInWhitelist) {
		return result, nil
	}
	if errors.Is(err, errorhandler.ErrUnauthenticated) {
		return result, status.Errorf(codes.Unauthenticated, "%v", err)
	}
	if errors.Is(err, errorhandler.ErrNoPermission) {
		return result, status.Errorf(codes.PermissionDenied, "%v", err)
	}
	if errors.Is(err, errorhandler.ErrInvalidArguments) {
		return result, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	if err != nil {
		return result, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	return result, nil
}
