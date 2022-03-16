package restful

import (
	"github.com/justdomepaul/toolbox/config"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type GinSuite struct {
	suite.Suite
	option        config.Set
	anotherOption config.Set
}

func (suite *GinSuite) SetupSuite() {
	suite.option = config.Set{
		Core: config.Core{
			SystemName: "Gin Server",
		},
		Server: config.Server{
			PrefixMessage:    "error gin server",
			AllowAllOrigins:  true,
			CustomizedRender: true,
		},
	}
	suite.anotherOption = config.Set{
		Core: config.Core{
			SystemName: "Gin Server",
		},
		Server: config.Server{
			ReleaseMode:     true,
			PrefixMessage:   "error gin server",
			AllowAllOrigins: false,
			AllowOrigins:    []string{"http://localhost"},
		},
	}
}

func (suite *GinSuite) TestNewGin() {
	gin, err := NewGin(suite.option, NewRender(), &JWTGuarder{})
	suite.NoError(err)
	suite.Equal("*gin.Engine", reflect.TypeOf(gin).String())
}

func (suite *GinSuite) TestNewGinAllowOrigins() {
	gin, err := NewGin(suite.option, NewRender(), &JWTGuarder{})
	suite.NoError(err)
	suite.Equal("*gin.Engine", reflect.TypeOf(gin).String())
}

func (suite *GinSuite) TestNewGinAllowOriginsReleaseAndLimitOrigin() {
	gin, err := NewGin(suite.anotherOption, NewRender(), &JWTGuarder{})
	suite.NoError(err)
	suite.Equal("*gin.Engine", reflect.TypeOf(gin).String())
}

func TestGinSuite(t *testing.T) {
	suite.Run(t, new(GinSuite))
}
