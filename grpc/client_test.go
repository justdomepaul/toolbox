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
	r, err := CreateClient("", config.GRPC{
		NoTLS:   true,
		SkipTLS: false,
	})
	suite.Equal("*grpc.ClientConn", reflect.TypeOf(r).String())
	suite.NoError(err)
}

func (suite *ClientSuite) TestCreateClientSKIPTLS() {
	r, err := CreateClient("", config.GRPC{
		NoTLS:   false,
		SkipTLS: true,
	})
	suite.Equal("*grpc.ClientConn", reflect.TypeOf(r).String())
	suite.NoError(err)
}

func (suite *ClientSuite) TestCreateClientALTS() {
	r, err := CreateClient("", config.GRPC{
		NoTLS:   false,
		SkipTLS: false,
		ALTS:    true,
	})
	suite.Equal("*grpc.ClientConn", reflect.TypeOf(r).String())
	suite.NoError(err)
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}
