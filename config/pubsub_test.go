package config

import (
	"github.com/stretchr/testify/suite"
	"os"
	"strconv"
	"testing"
)

type PubSubSuite struct {
	suite.Suite
	ProjectID             string
	EndPoint              string
	WithoutAuthentication bool
	GRPCInsecure          bool
}

func (suite *PubSubSuite) SetupSuite() {
	os.Clearenv()
	suite.ProjectID = "testProjectID"
	suite.EndPoint = "testEndPoint"
	suite.WithoutAuthentication = true
	suite.GRPCInsecure = true

	suite.NoError(os.Setenv("PROJECT_ID", suite.ProjectID))
	suite.NoError(os.Setenv("END_POINT", suite.EndPoint))
	suite.NoError(os.Setenv("WITHOUT_AUTHENTICATION", strconv.FormatBool(suite.WithoutAuthentication)))
	suite.NoError(os.Setenv("GRPC_INSECURE", strconv.FormatBool(suite.GRPCInsecure)))

}

func (suite *PubSubSuite) TestDefaultOption() {
	spanner := &PubSub{}
	suite.NoError(LoadFromEnv(spanner))
	suite.Equal(suite.ProjectID, spanner.ProjectID)
	suite.Equal(suite.EndPoint, spanner.EndPoint)
	suite.Equal(suite.WithoutAuthentication, spanner.WithoutAuthentication)
	suite.Equal(suite.GRPCInsecure, spanner.GRPCInsecure)
}

func TestPubSubSuite(t *testing.T) {
	suite.Run(t, new(PubSubSuite))
}
