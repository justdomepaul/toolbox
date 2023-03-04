package mongo

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/justdomepaul/toolbox/config"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"reflect"
	"testing"
)

type testISession struct {
	mock.Mock
	ISession
}

func (t *testISession) Close() error {
	args := t.Called()
	return args.Error(0)
}

type ConnectServiceSuite struct {
	suite.Suite
}

func (suite *ConnectServiceSuite) TestNewSession() {
	defer gostub.StubFunc(&newClient, &mongo.Client{}, nil).Reset()

	option := config.Mongo{}
	suite.NoError(config.LoadFromEnv(&option))
	_, errSession := NewSession(context.Background(), option)
	suite.NoError(errSession)
}

func (suite *ConnectServiceSuite) TestNewExtendMongoDatabase() {
	defer gostub.StubFunc(&NewSession, new(testISession), nil).Reset()

	option := config.Mongo{}
	suite.NoError(config.LoadFromEnv(&option))
	result, fn, err := NewExtendMongoDatabase(
		zap.NewExample(),
		option)
	suite.NoError(err)
	defer fn()
	suite.Equal("*mongo.testISession", reflect.TypeOf(result).String())
	suite.Equal("func()", reflect.TypeOf(fn).String())
}

func (suite *ConnectServiceSuite) TestNewExtendMongoDatabaseNewSessionError() {
	defer gostub.StubFunc(&NewSession, new(testISession), errors.New("got error")).Reset()
	option := config.Mongo{}
	suite.NoError(config.LoadFromEnv(&option))
	_, _, errExtendSpannerDatabase := NewExtendMongoDatabase(
		zap.NewExample(),
		option)
	suite.Error(errExtendSpannerDatabase)
}

func TestConnectServiceSuit(t *testing.T) {
	suite.Run(t, new(ConnectServiceSuite))
}
