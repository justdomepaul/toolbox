package grpc

import (
	"github.com/justdomepaul/toolbox/config"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
)

type ClientSuite struct {
	suite.Suite
}

func (suite *ClientSuite) SetupTest() {

}

func (suite *ClientSuite) TestCreateClientNO_TLS() {
	r, err := CreateClient("localhost:40800", config.GRPC{
		NoTLS:   true,
		SkipTLS: false,
	})
	suite.NoError(err)
	suite.Equal("*grpc.ClientConn", reflect.TypeOf(r).String())
}

func (suite *ClientSuite) TestCreateClientSKIPTLS() {
	r, err := CreateClient("localhost:40800", config.GRPC{
		NoTLS:   false,
		SkipTLS: true,
	})
	suite.NoError(err)
	suite.Equal("*grpc.ClientConn", reflect.TypeOf(r).String())
}

func (suite *ClientSuite) TestCreateClientALTS() {
	r, err := CreateClient("localhost:40800", config.GRPC{
		NoTLS:   false,
		SkipTLS: false,
		ALTS:    true,
	})
	suite.NoError(err)
	suite.Equal("*grpc.ClientConn", reflect.TypeOf(r).String())
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}
