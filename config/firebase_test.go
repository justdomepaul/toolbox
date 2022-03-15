package config

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type FirebaseSuite struct {
	suite.Suite
	FirebaseConfigJSON string
	FirebaseProjectID  string
}

func (suite *FirebaseSuite) SetupSuite() {
	os.Clearenv()
	suite.FirebaseConfigJSON = "testLinePayChannelID"
	suite.FirebaseProjectID = "testFirebaseProjectID"
	suite.NoError(os.Setenv("FIREBASE_CONFIG_JSON", suite.FirebaseConfigJSON))
	suite.NoError(os.Setenv("FIREBASE_PROJECT_ID", suite.FirebaseProjectID))
}

func (suite *FirebaseSuite) TestDefaultOption() {
	firebase := &Firebase{}
	suite.NoError(LoadFromEnv(firebase))
	suite.Equal(suite.FirebaseConfigJSON, firebase.FirebaseConfigJSON)
	suite.Equal(suite.FirebaseProjectID, firebase.FirebaseProjectID)
}

func TestFirebaseSuite(t *testing.T) {
	suite.Run(t, new(FirebaseSuite))
}
