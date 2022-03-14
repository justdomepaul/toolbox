package jwt

import (
	"github.com/cockroachdb/errors"
	"github.com/justdomepaul/toolbox/config"
	"github.com/prashantv/gostub"
	"github.com/square/go-jose/v3/jwt"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
	"time"
)

type JWTHS256Suite struct {
	suite.Suite
	key    string
	option config.JWT
}

func (suite *JWTHS256Suite) SetupTest() {
	suite.key = `b583ed184e2018b3d89a4fa8832d0a1f`
	result := config.JWT{}
	suite.NoError(config.LoadFromEnv(&result))
	suite.option = result
}

func (suite *JWTHS256Suite) TestNewHS256JWT() {
	suite.NotPanics(func() {
		_, err := NewHS256JWT(suite.key)
		suite.Equal(nil, err)
	})
}

func (suite *JWTHS256Suite) TestNewHS256JWTFromOptions() {
	suite.NotPanics(func() {
		_, err := NewHS256JWTFromOptions(suite.option)
		suite.NoError(err)
	})
}

func (suite *JWTHS256Suite) TestNewHS256JWTFromOptionsAllNoOption() {
	suite.NotPanics(func() {
		_, err := NewHS256JWTFromOptions(config.JWT{
			HmacSecretKeyPath: "",
			HmacSecretKey:     "",
		})
		suite.Error(err, ErrNoKey)
	})
}

func (suite *JWTHS256Suite) TestNewHS256JWTFromOptionsAllNoOptionNoFIle() {
	suite.NotPanics(func() {
		_, err := NewHS256JWTFromOptions(config.JWT{
			HmacSecretKeyPath: "testFile.txt",
			HmacSecretKey:     "",
		})
		suite.Error(err)
	})
}

func (suite *JWTHS256Suite) TestNewHS256JWTnewSignerError() {
	tSigner := testSigner{}
	defer gostub.StubFunc(&newSigner, tSigner, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewHS256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWTHS256Suite) TestHS256JWTGenerateTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewHS256JWT(suite.key)
		suite.Equal(nil, err)
		standClaims := NewClaimsBuilder().
			WithSubject("testTopic").
			WithIssuer("tester").
			WithID("test001").
			WithAudience([]string{"testerClient"}).
			ExpiresAfter(5 * time.Second).Build()
		//ExpiresAfter(87600 * time.Hour)
		common := NewCommon(standClaims, WithSecret("testData"))
		tk, errTk := j.GenerateToken(common)
		suite.T().Log(tk)
		suite.Equal(nil, errTk)
		suite.NotEmpty(tk)
	})
}

func (suite *JWTHS256Suite) TestHS256JWTValidateMethod() {
	suite.NotPanics(func() {
		j, err := NewHS256JWT(suite.key)
		suite.Equal(nil, err)

		suite.NoError(j.Validate(
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJQ2xhaW1zIjp7ImF1ZCI6InRlc3RlckNsaWVudCIsImV4cCI6MTYwNDMwODQwMCwiaXNzIjoidGVzdGVyIiwianRpIjoidGVzdDAwMSIsInN1YiI6InRlc3RUb3BpYyJ9LCJzIjoidGVzdERhdGEifQ.NLNmDo_Cq163gGEVFdQ_aFLNavnywlozw4XXhManHrE"))
	})
}

func (suite *JWTHS256Suite) TestHS256JWTValidateMethodparseSignedError() {
	suite.NotPanics(func() {
		defer gostub.StubFunc(&parseSigned, nil, errors.New("got error")).Reset()
		j, err := NewHS256JWT(suite.key)
		suite.Equal(nil, err)

		suite.Error(j.Validate(
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJQ2xhaW1zIjp7ImF1ZCI6InRlc3RlckNsaWVudCIsImV4cCI6MTYwNDMwODQwMCwiaXNzIjoidGVzdGVyIiwianRpIjoidGVzdDAwMSIsInN1YiI6InRlc3RUb3BpYyJ9LCJzIjoidGVzdERhdGEifQ.NLNmDo_Cq163gGEVFdQ_aFLNavnywlozw4XXhManHrE"))
	})
}

func (suite *JWTHS256Suite) TestHS256JWTParseTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewHS256JWT(suite.key)
		suite.Equal(nil, err)
		standClaims := NewClaimsBuilder().
			WithSubject("testTopic").
			WithIssuer("tester").
			WithID("test001").
			WithAudience([]string{"testerClient"}).
			ExpiresAfter(3 * time.Second).Build()
		common := NewCommon(standClaims, WithSecret("testData"))
		tk, errTk := j.GenerateToken(common)
		suite.Equal(nil, errTk)
		suite.T().Log(tk)
		standClaimsOutput := NewClaimsBuilder().Build()
		standCommonOutput := NewCommon(standClaimsOutput)
		errParse := j.VerifyToken(
			tk,
			standCommonOutput)
		suite.Equal(nil, errParse)
		suite.Equal("testTopic", standClaimsOutput.Subject)
		suite.Equal("tester", standClaimsOutput.Issuer)
		suite.Equal("test001", standClaimsOutput.ID)
		suite.Equal(jwt.Audience{"testerClient"}, standClaimsOutput.Audience)
	})
}

func (suite *JWTHS256Suite) TestHS256JWTParseTokenMethodExpire() {
	suite.NotPanics(func() {
		j, err := NewHS256JWT(suite.key)
		suite.Equal(nil, err)

		standCommonOutput := NewCommon(NewClaimsBuilder().Build())
		errParse := j.VerifyToken(
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0ZXJDbGllbnQiLCJleHAiOjE2MzMzMjAyMzcsImlzcyI6InRlc3RlciIsImp0aSI6InRlc3QwMDEiLCJzIjoidGVzdERhdGEiLCJzdWIiOiJ0ZXN0VG9waWMifQ.gR478g9oSrsQIL2uwQiShk38zvdxX3rtq9ofFQj2UDI",
			standCommonOutput)
		suite.Error(errParse, "token is expired")
	})
}

func (suite *JWTHS256Suite) TestHS256JWTParseTokenMethodExpireNoExpired() {
	suite.NotPanics(func() {
		j, err := NewHS256JWT(suite.key)
		suite.Equal(nil, err)

		standCommonOutput := NewCommon(NewClaimsBuilder().Build())
		errParse := j.VerifyToken(
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJQ2xhaW1zIjp7ImF1ZCI6InRlc3RlckNsaWVudCIsImlzcyI6InRlc3RlciIsImp0aSI6InRlc3QwMDEiLCJzdWIiOiJ0ZXN0VG9waWMifSwicyI6InRlc3REYXRhIn0.63AxNBLpjnYn5cER7MLrT7aISd3mfAasPtFcUE5Aup0",
			standCommonOutput)
		suite.NoError(errParse)
	})
}

func (suite *JWTHS256Suite) TestHS256JWTRefreshTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewHS256JWT(suite.key)
		suite.Equal(nil, err)
		standClaims := NewClaimsBuilder().
			WithSubject("testTopic").
			WithIssuer("tester").
			WithID("test001").
			WithAudience([]string{"testerClient"}).
			ExpiresAfter(3 * time.Second).Build()
		common := NewCommon(standClaims, WithSecret("testData"))
		tk, errTk := j.GenerateToken(common)
		log.Println(tk)
		suite.Equal(nil, errTk)
		standClaimsOutput := NewClaimsBuilder().Build()
		standCommonOutput := NewCommon(standClaimsOutput)
		newTk, errRefresh := j.RefreshToken(
			tk,
			standCommonOutput,
			100*time.Millisecond)
		suite.Equal(nil, errRefresh)
		suite.Equal(tk, newTk)
	})
}

func (suite *JWTHS256Suite) TestHS256JWTRefreshTokenMethodExpire() {
	suite.NotPanics(func() {
		j, err := NewHS256JWT(suite.key)
		suite.Equal(nil, err)
		standClaimsOutput := NewClaimsBuilder().Build()
		tk, errRefresh := j.RefreshToken(
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0ZXJDbGllbnQiLCJleHAiOjE2MzMzMjAyMzcsImlzcyI6InRlc3RlciIsImp0aSI6InRlc3QwMDEiLCJzIjoidGVzdERhdGEiLCJzdWIiOiJ0ZXN0VG9waWMifQ.gR478g9oSrsQIL2uwQiShk38zvdxX3rtq9ofFQj2UDI",
			standClaimsOutput,
			100*time.Millisecond,
		)
		suite.Equal(nil, errRefresh)
		suite.NotEmpty(tk)
		standClaimsOutputOther := NewClaimsBuilder().Build()
		standCommonOutput := NewCommon(standClaimsOutputOther)
		time.Sleep(105 * time.Millisecond)
		errParse := j.VerifyToken(tk, standCommonOutput)
		suite.Error(errParse, ErrTokenExpired)
	})
}

func (suite *JWTHS256Suite) TestHS256JWTRefreshTokenMethodparseRSRawError() {
	defer gostub.StubFunc(&parseSigned, nil, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		j, err := NewHS256JWT(suite.key)
		suite.Equal(nil, err)
		standClaimsOutput := NewClaimsBuilder().Build()
		_, errRefresh := j.RefreshToken(
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJQ2xhaW1zIjp7ImF1ZCI6InRlc3RlckNsaWVudCIsImV4cCI6MTYwNDMwODQwMCwiaXNzIjoidGVzdGVyIiwianRpIjoidGVzdDAwMSIsInN1YiI6InRlc3RUb3BpYyJ9LCJzIjoidGVzdERhdGEifQ.NLNmDo_Cq163gGEVFdQ_aFLNavnywlozw4XXhManHrE",
			standClaimsOutput,
			100*time.Millisecond,
		)
		suite.Error(errRefresh, "got error")
	})
}

func TestJWTHS256Suite(t *testing.T) {
	suite.Run(t, new(JWTHS256Suite))
}
