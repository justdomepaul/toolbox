package base58

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type Base58Suite struct {
	suite.Suite
	url string
}

func (suite *Base58Suite) SetupSuite() {
	suite.url = "https://www.amd.com/zh-hant/where-to-buy/radeon-rx-6000-series-graphics"
}

func (suite *Base58Suite) TestEncode() {
	result := Encode([]byte(suite.url))
	suite.Equal("MG5aEqKFcQiBEm1JraEWn22QMxkg39ejGatT2AKowLU6LqchF9daCwV88rQ6LUNmUZctckinnaELQHGn9QounF48dxTUN9cAE", result)
}

func (suite *Base58Suite) TestDecode() {
	result := Decode("MG5aEqKFcQiBEm1JraEWn22QMxkg39ejGatT2AKowLU6LqchF9daCwV88rQ6LUNmUZctckinnaELQHGn9QounF48dxTUN9cAE")
	suite.Equal([]byte{104, 116, 116, 112, 115, 58, 47, 47, 119, 119, 119, 46, 97, 109, 100, 46, 99, 111, 109, 47, 122, 104, 45, 104, 97, 110, 116, 47, 119, 104, 101, 114, 101, 45, 116, 111, 45, 98, 117, 121, 47, 114, 97, 100, 101, 111, 110, 45, 114, 120, 45, 54, 48, 48, 48, 45, 115, 101, 114, 105, 101, 115, 45, 103, 114, 97, 112, 104, 105, 99, 115}, result)
}

func TestBase58Suite(t *testing.T) {
	suite.Run(t, new(Base58Suite))
}
