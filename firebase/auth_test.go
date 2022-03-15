package firebase

import (
	"context"
	"github.com/justdomepaul/toolbox/config"
	"github.com/stretchr/testify/suite"
	"os"
	"reflect"
	"testing"
)

type AuthSuite struct {
	suite.Suite
	option config.Firebase
}

func (suite *AuthSuite) SetupSuite() {
	suite.NoError(os.Setenv("FIREBASE_CONFIG_JSON", `{
  "type": "service_account"
}`))
	option := config.Firebase{}
	suite.NoError(config.LoadFromEnv(&option))
	suite.option = option
}

func (suite *AuthSuite) TestNewAuthApp() {
	ctx := context.Background()
	client, err := NewClient(ctx, suite.option)
	suite.NoError(err)
	app, err := NewAuthApp(ctx, client)
	suite.NoError(err)
	suite.Equal("*auth.Client", reflect.TypeOf(app).String())
}

func (suite *AuthSuite) TestParseIDToken() {
	ctx := context.Background()
	client, err := NewClient(ctx, suite.option)
	suite.NoError(err)
	app, err := NewAuthApp(ctx, client)
	suite.NoError(err)
	token, err := VerifyIDToken(ctx, app, "testIDToken")
	suite.Error(err)
	suite.Nil(token)
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthSuite))
}
