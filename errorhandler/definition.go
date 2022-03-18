package errorhandler

import (
	"fmt"
	"github.com/cockroachdb/errors"
)

const (
	ErrProcessAuthenticate    = "errAuthenticate"
	ErrDbAlreadyExists        = "errDBAlreadyExists"
	ErrDbConnection           = "errDBConnection"
	ErrDbDisconnection        = "errDBDisConnection"
	ErrDbExecute              = "errDBExecute"
	ErrDbRowNotFound          = "errDBRowNotFound"
	ErrDbUpdateNoEffect       = "errDBUpdateNoEffect"
	ErrProcessExecute         = "errExecute"
	ErrGrpcConnection         = "errGRPCConnection"
	ErrGrpcExecute            = "errGRPCExecute"
	ErrProcessInvalidArgument = "errProcessInvalidArgument"
	ErrJsonMarshal            = "errJSONMarshal"
	ErrJsonUnmarshal          = "errJSONUnmarshal"
	ErrJwtExecute             = "errJWTExecute"
	ErrDataNotFound           = "errDataNotFound"
	ErrProcessPermissionDeny  = "errPermissionDeny"
	ErrProcessServerExecute   = "errServerExecute"
	ErrProcessVariable        = "errVariable"
)

var (
	ErrInWhitelist             = errors.New("in whitelist")
	ErrNoPermission            = errors.New("no permission")
	ErrUnauthenticated         = errors.New("unauthenticated")
	ErrInvalidArguments        = errors.New("invalid arguments")
	ErrNoRows                  = errors.New("no rows in result set")
	ErrAlreadyExists           = errors.New("primary key already exist")
	ErrUpdateNoEffect          = errors.New("no rows effected")
	ErrFailCloseSession        = errors.New("fail to close connection")
	ErrIncomingMetadataExist   = fmt.Errorf("%w: gRPC incoming metadata not exist", ErrInvalidArguments)
	ErrAuthorizationRequired   = fmt.Errorf("%w: metadata key: [authorization] must required", ErrInvalidArguments)
	ErrAuthorizationTypeBearer = fmt.Errorf("%w: JWT Authorization format error: must be Bearer", ErrInvalidArguments)
	ErrDeny                    = fmt.Errorf("%w: deny", ErrUnauthenticated)
	ErrNoRefreshToken          = fmt.Errorf("%w: no refresh token", ErrUnauthenticated)
	ErrExpiredToken            = fmt.Errorf("%w: token is expired", ErrUnauthenticated)
	ErrScopeNotExist           = fmt.Errorf("%w: scopes not exist", ErrInvalidArguments)
	ErrOutOfScopes             = fmt.Errorf("%w: out of scopes", ErrNoPermission)
	ErrOutOfPermissions        = fmt.Errorf("%w: out of permissions", ErrNoPermission)
)
