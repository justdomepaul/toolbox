package mongo

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/cockroachdb/errors"
	"github.com/justdomepaul/toolbox/config"
	"github.com/prashantv/gostub"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
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
	defer gostub.StubFunc(&newClient, &firestore.Client{}, nil).Reset()

	option := config.Firestore{}
	suite.NoError(config.LoadFromEnv(&option))
	_, errSession := NewSession(context.Background(), option)
	suite.NoError(errSession)
}

func (suite *ConnectServiceSuite) TestNewExtendFirestoreDatabase() {
	defer gostub.StubFunc(&NewSession, new(testISession), nil).Reset()

	option := config.Firestore{}
	suite.NoError(config.LoadFromEnv(&option))
	result, fn, err := NewExtendFirestoreDatabase(
		zap.NewExample(),
		option)
	suite.NoError(err)
	defer fn()
	suite.Equal("*firestore.testISession", reflect.TypeOf(result).String())
	suite.Equal("func()", reflect.TypeOf(fn).String())
}

func (suite *ConnectServiceSuite) TestNewExtendFirestoreDatabaseNewSessionError() {
	defer gostub.StubFunc(&NewSession, new(testISession), errors.New("got error")).Reset()
	option := config.Firestore{}
	suite.NoError(config.LoadFromEnv(&option))
	_, _, errExtendSpannerDatabase := NewExtendFirestoreDatabase(
		zap.NewExample(),
		option)
	suite.Error(errExtendSpannerDatabase)
}

func TestConnectServiceSuit(t *testing.T) {
	suite.Run(t, new(ConnectServiceSuite))
}
