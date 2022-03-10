package config

import (
	"github.com/stretchr/testify/assert"
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
	t := suite.T()
	os.Clearenv()
	suite.ProjectID = "testProjectID"
	suite.Instance = "testInstance"
	suite.Database = "testDatabase"
	suite.EndPoint = "testEndPoint"
	suite.WithoutAuthentication = true
	suite.GRPCInsecure = true

	assert.NoError(t, os.Setenv("PROJECT_ID", suite.ProjectID))
	assert.NoError(t, os.Setenv("INSTANCE", suite.Instance))
	assert.NoError(t, os.Setenv("DATABASE", suite.Database))
	assert.NoError(t, os.Setenv("END_POINT", suite.EndPoint))
	assert.NoError(t, os.Setenv("WITHOUT_AUTHENTICATION", strconv.FormatBool(suite.WithoutAuthentication)))
	assert.NoError(t, os.Setenv("GRPC_INSECURE", strconv.FormatBool(suite.GRPCInsecure)))

}

func (suite *SpannerSuite) TestDefaultOption() {
	t := suite.T()
	spanner := &Spanner{}
	suite.NoError(LoadFromEnv(spanner))
	assert.Equal(t, suite.ProjectID, spanner.ProjectID)
	assert.Equal(t, suite.Instance, spanner.Instance)
	assert.Equal(t, suite.Database, spanner.Database)
	assert.Equal(t, suite.EndPoint, spanner.EndPoint)
	assert.Equal(t, suite.WithoutAuthentication, spanner.WithoutAuthentication)
	assert.Equal(t, suite.GRPCInsecure, spanner.GRPCInsecure)
}

func TestSpannerSuite(t *testing.T) {
	suite.Run(t, new(SpannerSuite))
}
