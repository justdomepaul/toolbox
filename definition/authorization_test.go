package definition

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type AuthorizationSuite struct {
	suite.Suite
}

func (suite *AuthorizationSuite) TestGetConstant() {
	suite.Equal("authorization", AuthorizationKey)
	suite.Equal("authorization-client", AuthorizationClientKey)
	suite.Equal("Bearer ", AuthorizationType)
	suite.Equal("tk", QueryAuthKey)
	suite.Equal("tokenClaims", AuthTokenKey)
}

func TestAuthorizationSuite(t *testing.T) {
	suite.Run(t, new(AuthorizationSuite))
}
