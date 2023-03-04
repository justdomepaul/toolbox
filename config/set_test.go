package config

import (
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type ConfigSetSuite struct {
	suite.Suite
}

func (suite *ConfigSetSuite) TestNewConfigSet() {
	result, err := NewSet()
	suite.NoError(err)
	suite.Equal("config.Set", reflect.TypeOf(result).String())
}

func (suite *ConfigSetSuite) TestNewCassandra() {
	result, err := NewSet()
	suite.NoError(err)
	suite.Equal("Cassandra", reflect.TypeOf(NewCassandra(result)).Name())
}

func (suite *ConfigSetSuite) TestNewCloud() {
	result, err := NewSet()
	suite.NoError(err)
	suite.Equal("Cloud", reflect.TypeOf(NewCloud(result)).Name())
}

func (suite *ConfigSetSuite) TestNewCockroach() {
	result, err := NewSet()
	suite.NoError(err)
	suite.Equal("Cockroach", reflect.TypeOf(NewCockroach(result)).Name())
}

func (suite *ConfigSetSuite) TestNewCore() {
	result, err := NewSet()
	suite.NoError(err)
	suite.Equal("Core", reflect.TypeOf(NewCore(result)).Name())
}

func (suite *ConfigSetSuite) TestNewFirebase() {
	result, err := NewSet()
	suite.NoError(err)
	suite.Equal("Firebase", reflect.TypeOf(NewFirebase(result)).Name())
}

func (suite *ConfigSetSuite) TestNewGRPC() {
	result, err := NewSet()
	suite.NoError(err)
	suite.Equal("GRPC", reflect.TypeOf(NewGRPC(result)).Name())
}

func (suite *ConfigSetSuite) TestNewJWT() {
	result, err := NewSet()
	suite.NoError(err)
	suite.Equal("JWT", reflect.TypeOf(NewJWT(result)).Name())
}

func (suite *ConfigSetSuite) TestNewNewMongo() {
	result, err := NewSet()
	suite.NoError(err)
	suite.Equal("Mongo", reflect.TypeOf(NewMongo(result)).Name())
}

func (suite *ConfigSetSuite) TestNewPubSub() {
	result, err := NewSet()
	suite.NoError(err)
	suite.Equal("PubSub", reflect.TypeOf(NewPubSub(result)).Name())
}

func (suite *ConfigSetSuite) TestNewPostgres() {
	result, err := NewSet()
	suite.NoError(err)
	suite.Equal("Postgres", reflect.TypeOf(NewPostgres(result)).Name())
}

func (suite *ConfigSetSuite) TestNewRedis() {
	result, err := NewSet()
	suite.NoError(err)
	suite.Equal("Redis", reflect.TypeOf(NewRedis(result)).Name())
}

func (suite *ConfigSetSuite) TestNewServer() {
	result, err := NewSet()
	suite.NoError(err)
	suite.Equal("Server", reflect.TypeOf(NewServer(result)).Name())
}

func (suite *ConfigSetSuite) TestNewSpanner() {
	result, err := NewSet()
	suite.NoError(err)
	suite.Equal("Spanner", reflect.TypeOf(NewSpanner(result)).Name())
}

func TestConfigSetSuite(t *testing.T) {
	suite.Run(t, new(ConfigSetSuite))
}
