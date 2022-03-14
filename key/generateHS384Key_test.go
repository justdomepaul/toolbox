package key

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type GenerateHSKeySuite struct {
	suite.Suite
}

func (suite *GenerateHSKeySuite) TestGenerateHS256Key() {
	result, err := GenerateHS256Key()
	suite.NoError(err)
	suite.Len([]byte(result), 43)
	suite.T().Log(result)
	runes := ToBinaryRunes(result)
	suite.T().Log(runes)
	suite.Greater(len(runes), 256)
}

func (suite *GenerateHSKeySuite) TestGenerateHS384Key() {
	result, err := GenerateHS384Key()
	suite.NoError(err)
	suite.Len([]byte(result), 64)
	suite.T().Log(result)
	runes := ToBinaryRunes(result)
	suite.T().Log(runes)
	suite.Greater(len(runes), 384)
}

func (suite *GenerateHSKeySuite) TestGenerateHS512Key() {
	result, err := GenerateHS512Key()
	suite.NoError(err)
	suite.Len([]byte(result), 86)
	suite.T().Log(result)
	runes := ToBinaryRunes(result)
	suite.T().Log(runes)
	suite.Greater(len(runes), 512)
}

func TestGenerateHSKeySuite(t *testing.T) {
	suite.Run(t, new(GenerateHSKeySuite))
}
