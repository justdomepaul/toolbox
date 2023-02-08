package config

import (
	"github.com/stretchr/testify/suite"
	"os"
	"strconv"
	"testing"
)

type RedisSuite struct {
	suite.Suite
	RedisUsername   string
	RedisPassword   string
	RedisClientName string
	RedisHost       string
	RedisPort       string
	RedisPoolSize   int
}

func (suite *RedisSuite) SetupSuite() {
	os.Clearenv()
	suite.RedisUsername = "testRedisUsername"
	suite.RedisPassword = "testRedisPassword"
	suite.RedisClientName = "testRedisClientName"
	suite.RedisHost = "testRedisHost"
	suite.RedisPort = "testRedisPort"
	suite.RedisPoolSize = 100

	suite.NoError(os.Setenv("REDIS_USERNAME", suite.RedisUsername))
	suite.NoError(os.Setenv("REDIS_PASSWORD", suite.RedisPassword))
	suite.NoError(os.Setenv("REDIS_CLIENT_NAME", suite.RedisClientName))
	suite.NoError(os.Setenv("REDIS_HOST", suite.RedisHost))
	suite.NoError(os.Setenv("REDIS_PORT", suite.RedisPort))
	suite.NoError(os.Setenv("REDIS_POOL_SIZE", strconv.FormatInt(int64(suite.RedisPoolSize), 10)))

}

func (suite *RedisSuite) TestDefaultOption() {
	redis := &Redis{}
	suite.NoError(LoadFromEnv(redis))
	suite.Equal(suite.RedisUsername, redis.RedisUsername)
	suite.Equal(suite.RedisPassword, redis.RedisPassword)
	suite.Equal(suite.RedisClientName, redis.RedisClientName)
	suite.Equal(suite.RedisHost, redis.RedisHost)
	suite.Equal(suite.RedisPort, redis.RedisPort)
	suite.Equal(suite.RedisPoolSize, redis.RedisPoolSize)
}

func TestRedisSuite(t *testing.T) {
	suite.Run(t, new(RedisSuite))
}
