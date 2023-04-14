package config

import (
	"github.com/stretchr/testify/suite"
	"os"
	"strconv"
	"testing"
)

type MongoSuite struct {
	suite.Suite
	MongoProtocol          string
	MongoUsername          string
	MongoPassword          string
	MongoHost              string
	MongoDatabase          string
	MongoAuthSource        bool
	MongoIgnoreQueryString bool
}

func (suite *MongoSuite) SetupSuite() {
	os.Clearenv()
	suite.MongoProtocol = "testMongoProtocol"
	suite.MongoUsername = "testMongoUsername"
	suite.MongoPassword = "testMongoPassword"
	suite.MongoHost = "testMongoHost"
	suite.MongoDatabase = "testMongoDatabase"
	suite.MongoAuthSource = false
	suite.MongoIgnoreQueryString = false

	suite.NoError(os.Setenv("MONGO_PROTOCOL", suite.MongoProtocol))
	suite.NoError(os.Setenv("MONGO_USERNAME", suite.MongoUsername))
	suite.NoError(os.Setenv("MONGO_PASSWORD", suite.MongoPassword))
	suite.NoError(os.Setenv("MONGO_HOST", suite.MongoHost))
	suite.NoError(os.Setenv("MONGO_DATABASE", suite.MongoDatabase))
	suite.NoError(os.Setenv("MONGO_AUTH_SOURCE", strconv.FormatBool(suite.MongoAuthSource)))
	suite.NoError(os.Setenv("MONGO_IGNORE_QUERY_STRING", strconv.FormatBool(suite.MongoIgnoreQueryString)))
}

func (suite *MongoSuite) TestDefaultOption() {
	mongo := &Mongo{}
	suite.NoError(LoadFromEnv(mongo))
	suite.Equal(suite.MongoProtocol, mongo.MongoProtocol)
	suite.Equal(suite.MongoUsername, mongo.MongoUsername)
	suite.Equal(suite.MongoPassword, mongo.MongoPassword)
	suite.Equal(suite.MongoHost, mongo.MongoHost)
	suite.Equal(suite.MongoDatabase, mongo.MongoDatabase)
	suite.Equal(suite.MongoAuthSource, mongo.MongoAuthSource)
	suite.Equal(suite.MongoIgnoreQueryString, mongo.MongoIgnoreQueryString)
}

func TestMongoSuite(t *testing.T) {
	suite.Run(t, new(MongoSuite))
}
