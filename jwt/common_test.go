package jwt

import (
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
	"time"
)

type CommonSuite struct {
	suite.Suite
}

func (suite *CommonSuite) TestNewCommon() {
	tk := NewCommon(
		NewClaimsBuilder().WithSubject("testSubject").WithIssuer("testIssuer").ExpiresAfter(100*time.Second).Build(),
		WithSecret("testSecret"),
	)

	suite.Equal("*jwt.Common", reflect.TypeOf(tk).String())
	suite.Equal("testSecret", tk.Secret)
	result, err := json.Marshal(tk)
	suite.NoError(err)
	suite.T().Log(string(result))
}

func TestCommonSuite(t *testing.T) {
	suite.Run(t, new(CommonSuite))
}
