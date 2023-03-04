package config

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type MongoSuite struct {
	suite.Suite
	MongoProtocol string
	MongoUsername string
	MongoPassword string
	MongoHost     string
	MongoDatabase string
}

func (suite *MongoSuite) SetupSuite() {
	os.Clearenv()
	suite.MongoProtocol = "testMongoProtocol"
	suite.MongoUsername = "testMongoUsername"
	suite.MongoPassword = "testMongoPassword"
	suite.MongoHost = "testMongoHost"
	suite.MongoDatabase = "testMongoDatabase"

	suite.NoError(os.Setenv("MONGO_PROTOCOL", suite.MongoProtocol))
	suite.NoError(os.Setenv("MONGO_USERNAME", suite.MongoUsername))
	suite.NoError(os.Setenv("MONGO_PASSWORD", suite.MongoPassword))
	suite.NoError(os.Setenv("MONGO_HOST", suite.MongoHost))
	suite.NoError(os.Setenv("MONGO_DATABASE", suite.MongoDatabase))
}

func (suite *MongoSuite) TestDefaultOption() {
	mongo := &Mongo{}
	suite.NoError(LoadFromEnv(mongo))
	suite.Equal(suite.MongoProtocol, mongo.MongoProtocol)
	suite.Equal(suite.MongoUsername, mongo.MongoUsername)
	suite.Equal(suite.MongoPassword, mongo.MongoPassword)
	suite.Equal(suite.MongoHost, mongo.MongoHost)
	suite.Equal(suite.MongoDatabase, mongo.MongoDatabase)
}

func TestMongoSuite(t *testing.T) {
	suite.Run(t, new(MongoSuite))
}
