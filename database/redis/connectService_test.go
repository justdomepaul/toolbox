package redis

import (
	"github.com/cockroachdb/errors"
	"github.com/justdomepaul/toolbox/config"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

type ConnectServiceSuite struct {
	suite.Suite
}

func (suite *ConnectServiceSuite) TestNewCluster() {
	client, err := NewSession(config.Redis{
		RedisUsername:   "root",
		RedisPassword:   "root",
		RedisClientName: "root",
		RedisPoolSize:   11,
	})
	suite.NoError(err)
	suite.Equal("*redis.Client", reflect.TypeOf(client).String())
}

func (suite *ConnectServiceSuite) TestNewExtendRedisDatabase() {
	result, fn, err := NewExtendRedisDatabase(
		zap.NewExample(),
		config.Redis{
			RedisUsername:   "root",
			RedisPassword:   "root",
			RedisClientName: "root",
			RedisPoolSize:   11,
		})
	suite.NoError(err)
	suite.Equal("*redis.Client", reflect.TypeOf(result).String())
	suite.Equal("func()", reflect.TypeOf(fn).String())
	suite.NotPanics(func() {
		fn()
	})
}

func (suite *ConnectServiceSuite) TestNewExtendRedisDatabaseNewSessionError() {
	defer gostub.StubFunc(&NewSession, nil, errors.New("got error")).Reset()
	_, _, err := NewExtendRedisDatabase(
		zap.NewExample(),
		config.Redis{})
	suite.Error(err)
}

func TestConnectServiceSuite(t *testing.T) {
	suite.Run(t, new(ConnectServiceSuite))
}
