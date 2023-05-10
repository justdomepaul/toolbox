package firebase

import (
	"context"
	"firebase.google.com/go/messaging"
	"github.com/justdomepaul/toolbox/config"
	"github.com/stretchr/testify/suite"
	"os"
	"reflect"
	"testing"
)

type CloudMessaginggSuite struct {
	suite.Suite
	option config.Firebase
}

func (suite *CloudMessaginggSuite) SetupSuite() {
	suite.NoError(os.Setenv("FIREBASE_CONFIG_JSON_BASE64", `ewogICJ0eXBlIjogInNlcnZpY2VfYWNjb3VudCIKfQ==`))
	option := config.Firebase{}
	suite.NoError(config.LoadFromEnv(&option))
	option.FirebaseConfigJSON = ""
	option.FirebaseProjectID = "toolbox-8d1fc"
	suite.option = option
}

func (suite *CloudMessaginggSuite) TestNewCloudMessagingApp() {
	ctx := context.Background()
	client, err := NewClient(ctx, suite.option)
	suite.NoError(err)
	app, err := NewCloudMessagingApp(ctx, client)
	suite.NoError(err)
	app.Send(ctx, &messaging.Message{
		Data:  nil,
		Topic: "",
	})
	suite.Equal("*messaging.Client", reflect.TypeOf(app).String())
}

func TestCloudMessagingSuite(t *testing.T) {
	suite.Run(t, new(CloudMessaginggSuite))
}
