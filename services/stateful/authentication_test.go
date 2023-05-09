package stateful

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/justdomepaul/toolbox/config"
	"github.com/justdomepaul/toolbox/errorhandler"
	"github.com/justdomepaul/toolbox/jwt"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type CommonAuthenticationSuite struct {
	suite.Suite
	uid    uuid.UUID
	id     []byte
	key    string
	option config.JWT
	jwt    jwt.IJWT
	token  string
}

func (suite *CommonAuthenticationSuite) SetupSuite() {
	uid, err := uuid.NewRandom()
	suite.NoError(err)
	suite.uid = uid
	suite.id = uid[:]

	suite.key = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIANpsaiXN0UdToO728+5kwGNPA+RsbEPwb3MGiAlBhSXoAoGCCqGSM49
AwEHoUQDQgAER5WNvPs/SMICGESgDbN7IYl0CvPSkUhAaUtF/LAQEINqte/HLMks
hRsKJ2MTCe1upn5vhgBuGl5CL4ea4DqNhA==
-----END EC PRIVATE KEY-----
`
	result := config.JWT{}
	suite.NoError(config.LoadFromEnv(&result))
	suite.option = result

	j, err := jwt.NewES256JWT(suite.key)
	suite.NoError(err)
	suite.jwt = j

	tk := newToken(jwt.NewClaimsBuilder().Build(), jwt.WithClientID(suite.id), jwt.WithScopes("/ping", "/pong"))
	token, err := j.GenerateToken(tk)
	suite.NoError(err)
	suite.token = token
}

func (suite *CommonAuthenticationSuite) TestNewCommonAuthentication() {
	gOp := config.GRPC{}
	suite.NoError(config.LoadFromEnv(&gOp))
	service, err := NewAuthentication(gOp, suite.jwt)
	suite.NoError(err)
	suite.Equal("*stateful.Authentication", reflect.TypeOf(service).String())
}

func (suite *CommonAuthenticationSuite) TestAuthenticateMethod() {
	ctx := context.Background()
	gOp := config.GRPC{
		AllowedList: []string{
			"/ping",
		},
	}
	suite.NoError(config.LoadFromEnv(&gOp))
	service, err := NewAuthentication(gOp, suite.jwt)
	suite.NoError(err)
	result, err := service.Authenticate(ctx, func() (string, error) {
		return suite.token, nil
	}, "/pong")
	suite.NoError(err)
	suite.Equal(suite.id, result.GetID())
}

func (suite *CommonAuthenticationSuite) TestAuthenticateMethodInWhiteList() {
	ctx := context.Background()
	gOp := config.GRPC{
		AllowedList: []string{
			"/ping",
		},
	}
	service, err := NewAuthentication(gOp, suite.jwt)
	suite.NoError(err)
	resultID, err := service.Authenticate(ctx, func() (string, error) {
		return suite.token, nil
	}, "/ping")
	suite.ErrorIs(err, errorhandler.ErrInWhitelist)
	suite.Empty(resultID)
}

func (suite *CommonAuthenticationSuite) TestAuthenticateMethodTokenFnError() {
	ctx := context.Background()
	gOp := config.GRPC{
		AllowedList: []string{
			"/ping",
		},
	}
	service, err := NewAuthentication(gOp, suite.jwt)
	suite.NoError(err)
	resultID, err := service.Authenticate(ctx, func() (string, error) {
		return suite.token, errors.New("got error")
	}, "/pong")
	suite.ErrorIs(err, errorhandler.ErrUnauthenticated)
	suite.Empty(resultID)
}

func (suite *CommonAuthenticationSuite) TestAuthenticateMethodVerifyTokenError() {
	ctx := context.Background()
	gOp := config.GRPC{
		AllowedList: []string{
			"/ping",
		},
	}
	service, err := NewAuthentication(gOp, suite.jwt)
	suite.NoError(err)
	resultID, err := service.Authenticate(ctx, func() (string, error) {
		return "test token", nil
	}, "/pong")
	suite.ErrorIs(err, errorhandler.ErrUnauthenticated)
	suite.Empty(resultID)
}

func (suite *CommonAuthenticationSuite) TestAuthenticateMethodOutOfScopes() {
	ctx := context.Background()
	gOp := config.GRPC{
		AllowedList: []string{
			"/ping",
		},
	}
	suite.NoError(config.LoadFromEnv(&gOp))
	service, err := NewAuthentication(gOp, suite.jwt)
	suite.NoError(err)
	resultID, err := service.Authenticate(ctx, func() (string, error) {
		return suite.token, nil
	}, "/foo")
	suite.ErrorIs(err, errorhandler.ErrOutOfScopes)
	suite.Empty(resultID)
}

func TestCommonAuthenticationSuite(t *testing.T) {
	suite.Run(t, new(CommonAuthenticationSuite))
}
