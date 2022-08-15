package firebase

import (
	"context"
	"github.com/justdomepaul/toolbox/config"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type FirestoreSuite struct {
	suite.Suite
	option config.Firebase
}

func (suite *FirestoreSuite) SetupSuite() {
	option := config.Firebase{}
	suite.NoError(config.LoadFromEnv(&option))
	option.FirebaseConfigJSON = ""
	option.FirebaseProjectID = "toolbox-8d1fc"
	suite.option = option
}

func (suite *FirestoreSuite) TestNewFirestoreApp() {
	ctx := context.Background()
	client, err := NewClient(ctx, suite.option)
	suite.NoError(err)
	app, err := NewFirestoreApp(ctx, client)
	suite.NoError(err)
	suite.Equal("*firestore.Client", reflect.TypeOf(app).String())
}

func TestFirestoreSuite(t *testing.T) {
	suite.Run(t, new(FirestoreSuite))
}
