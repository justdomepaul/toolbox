package utils

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

type UUIDSuite struct {
	suite.Suite
}

func (suite *UUIDSuite) TestParseUUID() {
	uid, err := uuid.NewRandom()
	suite.NoError(err)
	suite.NotPanics(func() {
		suite.Equal(uid[:], ParseUUID(uid.String()))
	})
}

func (suite *UUIDSuite) TestParseUUIDEmpty() {
	suite.NotPanics(func() {
		suite.Empty(ParseUUID(""))
	})
}

func (suite *UUIDSuite) TestParseUUIDError() {
	suite.Panics(func() {
		ParseUUID("testUUID")
	})
}

func (suite *UUIDSuite) TestFromUUID() {
	uid, err := uuid.NewRandom()
	suite.NoError(err)
	suite.NotPanics(func() {
		suite.Equal(uid.String(), FromUUID(uid[:]))
	})
}

func (suite *UUIDSuite) TestFromUUIDEmpty() {
	suite.NotPanics(func() {
		suite.Equal("", FromUUID([]byte(nil)))
	})
}

func (suite *UUIDSuite) TestFromUUIDError() {
	suite.Panics(func() {
		FromUUID([]byte("testUUID"))
	})
}

func TestUUIDSuite(t *testing.T) {
	suite.Run(t, new(UUIDSuite))
}
