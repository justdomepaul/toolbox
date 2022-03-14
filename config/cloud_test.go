package config

import (
	"github.com/stretchr/testify/suite"
	"os"
	"strconv"
	"testing"
)

type CloudSuite struct {
	suite.Suite
	ProjectID             string
	EndPoint              string
	WithoutAuthentication bool
	GRPCInsecure          bool
}

func (suite *CloudSuite) SetupSuite() {
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

func (suite *CloudSuite) TestDefaultOption() {
	spanner := &Cloud{}
	suite.NoError(LoadFromEnv(spanner))
	suite.Equal(suite.ProjectID, spanner.ProjectID)
	suite.Equal(suite.EndPoint, spanner.EndPoint)
	suite.Equal(suite.WithoutAuthentication, spanner.WithoutAuthentication)
	suite.Equal(suite.GRPCInsecure, spanner.GRPCInsecure)
}

func TestCloudSuite(t *testing.T) {
	suite.Run(t, new(CloudSuite))
}
