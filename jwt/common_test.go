package jwt

import (
	"encoding/json"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
	"time"
)

// NewMockClaim method
func NewMockClaim(claims *jwt.Claims) *mockClaim {
	token := &mockClaim{
		Claims: claims,
	}
	return token
}

type mockClaim struct {
	*jwt.Claims
}

type CommonSuite struct {
	suite.Suite
}

func (suite *CommonSuite) TestNewCommon() {
	tk := NewCommon(
		NewClaimsBuilder().WithSubject("testSubject").WithIssuer("testIssuer").ExpiresAfter(100*time.Second).Build(),
		WithSecret("testSecret"),
		WithPermissions("/ping", "/pong"),
		WithScopes("/ping", "/pong"),
	)

	suite.Equal("*jwt.Common", reflect.TypeOf(tk).String())
	suite.Equal([]byte("testSecret"), tk.Secret)
	result, err := json.Marshal(tk)
	suite.NoError(err)
	suite.T().Log(string(result))
}

func TestCommonSuite(t *testing.T) {
	suite.Run(t, new(CommonSuite))
}
