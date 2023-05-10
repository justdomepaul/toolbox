package firebase

import (
	"context"
	"github.com/justdomepaul/toolbox/config"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type ClientSuite struct {
	suite.Suite
	option config.Firebase
}

func (suite *ClientSuite) SetupSuite() {
	suite.NoError(os.Setenv("FIREBASE_CONFIG_JSON", `{}`))
	option := config.Firebase{}
	suite.NoError(config.LoadFromEnv(&option))
	suite.option = option
}

func (suite *ClientSuite) TestNewApp() {
	app, errClient := NewClient(context.Background(), config.Firebase{
		FirebaseConfigJSON:       "",
		FirebaseConfigJSONBase64: "",
		FirebaseProjectID:        "microbee",
	})
	suite.NoError(errClient)
	suite.T().Log(app)
}

func (suite *ClientSuite) TestNewAppHaveJson() {
	app, errClient := NewClient(context.Background(), suite.option)
	suite.NoError(errClient)
	suite.T().Log(app)
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}
