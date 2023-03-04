package config

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type MongoSuite struct {
	suite.Suite
	MongoUsername string
	MongoPassword string
	MongoHost     string
	MongoDatabase string
}

func (suite *MongoSuite) SetupSuite() {
	os.Clearenv()
	suite.MongoUsername = "testMongoUsername"
	suite.MongoPassword = "testMongoPassword"
	suite.MongoHost = "testMongoHost"
	suite.MongoDatabase = "testMongoDatabase"

	suite.NoError(os.Setenv("MONGO_USERNAME", suite.MongoUsername))
	suite.NoError(os.Setenv("MONGO_PASSWORD", suite.MongoPassword))
	suite.NoError(os.Setenv("MONGO_HOST", suite.MongoHost))
	suite.NoError(os.Setenv("MONGO_DATABASE", suite.MongoDatabase))
}

func (suite *MongoSuite) TestDefaultOption() {
	firebase := &Mongo{}
	suite.NoError(LoadFromEnv(firebase))
	suite.Equal(suite.MongoUsername, firebase.MongoUsername)
	suite.Equal(suite.MongoPassword, firebase.MongoPassword)
	suite.Equal(suite.MongoHost, firebase.MongoHost)
	suite.Equal(suite.MongoDatabase, firebase.MongoDatabase)
}

func TestMongoSuite(t *testing.T) {
	suite.Run(t, new(MongoSuite))
}
