package gormspanner

import (
	"context"
	"github.com/justdomepaul/toolbox/config"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gorm.io/gorm/schema"
	"reflect"
	"testing"
)

type CamelCaseNamerSuite struct {
	suite.Suite
}

func (suite *CamelCaseNamerSuite) TestTableName() {
	suite.Equal(NewCamelCaseNamer().TableName("testTableName"), "testTableName")
}

func (suite *CamelCaseNamerSuite) TestSchemaName() {
	suite.Equal(NewCamelCaseNamer().SchemaName("testSchemaName"), "testSchemaName")
}

func (suite *CamelCaseNamerSuite) TestColumnName() {
	suite.Equal(NewCamelCaseNamer().ColumnName("", "testColumnName"), "testColumnName")
}

func (suite *CamelCaseNamerSuite) TestJoinTableName() {
	suite.Equal(NewCamelCaseNamer().JoinTableName("testJoinTableName"), "testJoinTableName")
}

func (suite *CamelCaseNamerSuite) TestRelationshipFKName() {
	suite.Equal(NewCamelCaseNamer().RelationshipFKName(schema.Relationship{
		Name: "testRelationshipFKName",
	}), "testRelationshipFKName")
}

func (suite *CamelCaseNamerSuite) TestCheckerName() {
	suite.Equal(NewCamelCaseNamer().CheckerName("testTableName", "testCheckerName"), "testTableName_testCheckerName")
}

func (suite *CamelCaseNamerSuite) TestIndexName() {
	suite.Equal(NewCamelCaseNamer().IndexName("", "testIndexName"), "testIndexNameIdx")
}

func (suite *CamelCaseNamerSuite) TestUniqueName() {
	suite.Equal(NewCamelCaseNamer().UniqueName("", "testUniqueName"), "testUniqueNameUIdx")
}

func TestCamelCaseNamerSuite(t *testing.T) {
	suite.Run(t, new(CamelCaseNamerSuite))
}

type ConnectServiceSuite struct {
	suite.Suite
}

func (suite *ConnectServiceSuite) TestNewSession() {
	option := config.Spanner{}
	suite.NoError(config.LoadFromEnv(&option))
	_, errSession := NewSession(context.Background(), option)
	suite.NoError(errSession)
}

func (suite *ConnectServiceSuite) TestNewExtendGormSpannerDatabase() {
	option := config.Spanner{}
	suite.NoError(config.LoadFromEnv(&option))
	result, fn, err := NewExtendGormSpannerDatabase(
		zap.NewExample(),
		option)
	suite.NoError(err)
	defer fn()
	suite.Equal("*gorm.DB", reflect.TypeOf(result).String())
	suite.Equal("func()", reflect.TypeOf(fn).String())
}

func TestConnectServiceSuite(t *testing.T) {
	suite.Run(t, new(ConnectServiceSuite))
}
