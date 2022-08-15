package config

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type FirestoreSuite struct {
	suite.Suite
	ProjectID string
	EndPoint  string
}

func (suite *FirestoreSuite) SetupSuite() {
	os.Clearenv()
	suite.ProjectID = "testProjectID"
	suite.EndPoint = "testEndPoint"

	suite.NoError(os.Setenv("PROJECT_ID", suite.ProjectID))
	suite.NoError(os.Setenv("END_POINT", suite.EndPoint))
}

func (suite *FirestoreSuite) TestDefaultOption() {
	firebase := &Firestore{}
	suite.NoError(LoadFromEnv(firebase))
	suite.Equal(suite.ProjectID, firebase.ProjectID)
	suite.Equal(suite.EndPoint, firebase.EndPoint)
}

func TestFirestoreSuite(t *testing.T) {
	suite.Run(t, new(FirestoreSuite))
}
