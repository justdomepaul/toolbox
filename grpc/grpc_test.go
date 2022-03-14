package grpc

import (
	"github.com/justdomepaul/toolbox/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type ClientSuit struct {
	suite.Suite
}

func (suite *ClientSuit) SetupTest() {

}

func (suite *ClientSuit) TestCreateClientNO_TLS() {
	t := suite.T()
	r, err := CreateClient("", config.GRPC{
		NoTLS:   true,
		SkipTLS: false,
	})
	assert.Equal(t, "*grpc.ClientConn", reflect.TypeOf(r).String())
	assert.NoError(t, err)
}

func (suite *ClientSuit) TestCreateClientSKIP_TLS() {
	t := suite.T()
	r, err := CreateClient("", config.GRPC{
		NoTLS:   false,
		SkipTLS: true,
	})
	assert.Equal(t, "*grpc.ClientConn", reflect.TypeOf(r).String())
	assert.NoError(t, err)
}

func TestClientSuit(t *testing.T) {
	suite.Run(t, new(ClientSuit))
}
