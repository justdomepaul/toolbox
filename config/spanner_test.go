package config

import (
	"github.com/stretchr/testify/suite"
	"os"
	"strconv"
	"testing"
)

type SpannerSuite struct {
	suite.Suite
	ProjectID             string
	Instance              string
	Database              string
	EndPoint              string
	WithoutAuthentication bool
	GRPCInsecure          bool
}

func (suite *SpannerSuite) SetupSuite() {
	os.Clearenv()
	suite.ProjectID = "testProjectID"
	suite.Instance = "testInstance"
	suite.Database = "testDatabase"
	suite.EndPoint = "testEndPoint"
	suite.WithoutAuthentication = true
	suite.GRPCInsecure = true

	suite.NoError(os.Setenv("PROJECT_ID", suite.ProjectID))
	suite.NoError(os.Setenv("INSTANCE", suite.Instance))
	suite.NoError(os.Setenv("DATABASE", suite.Database))
	suite.NoError(os.Setenv("END_POINT", suite.EndPoint))
	suite.NoError(os.Setenv("WITHOUT_AUTHENTICATION", strconv.FormatBool(suite.WithoutAuthentication)))
	suite.NoError(os.Setenv("GRPC_INSECURE", strconv.FormatBool(suite.GRPCInsecure)))

}

func (suite *SpannerSuite) TestDefaultOption() {
	spanner := &Spanner{}
	suite.NoError(LoadFromEnv(spanner))
	suite.Equal(suite.ProjectID, spanner.ProjectID)
	suite.Equal(suite.Instance, spanner.Instance)
	suite.Equal(suite.Database, spanner.Database)
	suite.Equal(suite.EndPoint, spanner.EndPoint)
	suite.Equal(suite.WithoutAuthentication, spanner.WithoutAuthentication)
	suite.Equal(suite.GRPCInsecure, spanner.GRPCInsecure)
}

func TestSpannerSuite(t *testing.T) {
	suite.Run(t, new(SpannerSuite))
}
