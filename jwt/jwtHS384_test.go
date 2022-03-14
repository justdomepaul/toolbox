package jwt

import (
	"github.com/cockroachdb/errors"
	"github.com/justdomepaul/toolbox/config"
	"github.com/justdomepaul/toolbox/key"
	"github.com/prashantv/gostub"
	"github.com/square/go-jose/v3/jwt"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
	"time"
)

type testClaims struct {
	Email string `json:"email,omitempty"`
	*jwt.Claims
}

type JWTHS384Suite struct {
	suite.Suite
	key    string
	option config.JWT
}

func (suite *JWTHS384Suite) SetupTest() {
	suite.key = `b583ed184e2018b3d89a4fa8832d0a1f`
	result := config.JWT{}
	suite.NoError(config.LoadFromEnv(&result))
	suite.option = result
}

func (suite *JWTHS384Suite) TestNewHS384JWT() {
	suite.NotPanics(func() {
		_, err := NewHS384JWT(suite.key)
		suite.Equal(nil, err)
	})
}

func (suite *JWTHS384Suite) TestNewHS384JWTByHS384Key() {
	testEmail := "testMock@mock.com"
	key384, err := key.GenerateHS384Key()
	suite.NoError(err)
	hs384JWT, err := NewHS384JWT(key384)
	suite.NoError(err)
	standClaims := NewClaimsBuilder().
		WithSubject("testTopic").
		WithIssuer("tester").
		WithID("test001").
		WithAudience([]string{"testerClient"}).
		ExpiresAfter(5 * time.Second).Build()
	//ExpiresAfter(87600 * time.Hour)
	testClaimsInput := &testClaims{
		Claims: standClaims,
		Email:  testEmail,
	}
	result, err := hs384JWT.GenerateToken(testClaimsInput)
	suite.NoError(err)
	suite.T().Log(result)
	testClaimsResult := &testClaims{
		Claims: NewClaimsBuilder().Build(),
	}
	suite.NoError(ParseUnverified(result, testClaimsResult))
	suite.Equal(testEmail, testClaimsResult.Email)
}

func (suite *JWTHS384Suite) TestNewHS384JWTFromOptions() {
	suite.NotPanics(func() {
		_, err := NewHS384JWTFromOptions(suite.option)
		suite.NoError(err)
	})
}

func (suite *JWTHS384Suite) TestNewHS384JWTFromOptionsAllNoOption() {
	suite.NotPanics(func() {
		_, err := NewHS384JWTFromOptions(config.JWT{
			HmacSecretKeyPath: "",
			HmacSecretKey:     "",
		})
		suite.Error(err, ErrNoKey)
	})
}

func (suite *JWTHS384Suite) TestNewHS384JWTFromOptionsAllNoOptionNoFIle() {
	suite.NotPanics(func() {
		_, err := NewHS384JWTFromOptions(config.JWT{
			HmacSecretKeyPath: "testFile.txt",
			HmacSecretKey:     "",
		})
		suite.Error(err)
	})
}

func (suite *JWTHS384Suite) TestNewHS384JWTnewSignerError() {
	tSigner := testSigner{}
	defer gostub.StubFunc(&newSigner, tSigner, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewHS384JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWTHS384Suite) TestHS384JWTGenerateTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewHS384JWT(suite.key)
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

func (suite *JWTHS384Suite) TestHS384JWTValidateMethod() {
	suite.NotPanics(func() {
		j, err := NewHS384JWT(suite.key)
		suite.Equal(nil, err)

		suite.NoError(j.Validate(
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJQ2xhaW1zIjp7ImF1ZCI6InRlc3RlckNsaWVudCIsImV4cCI6MTYwNDMwODQwMCwiaXNzIjoidGVzdGVyIiwianRpIjoidGVzdDAwMSIsInN1YiI6InRlc3RUb3BpYyJ9LCJzIjoidGVzdERhdGEifQ.NLNmDo_Cq163gGEVFdQ_aFLNavnywlozw4XXhManHrE"))
	})
}

func (suite *JWTHS384Suite) TestHS384JWTValidateMethodparseSignedError() {
	suite.NotPanics(func() {
		defer gostub.StubFunc(&parseSigned, nil, errors.New("got error")).Reset()
		j, err := NewHS384JWT(suite.key)
		suite.Equal(nil, err)

		suite.Error(j.Validate(
			"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJQ2xhaW1zIjp7ImF1ZCI6InRlc3RlckNsaWVudCIsImV4cCI6MTYwNDMwODQwMCwiaXNzIjoidGVzdGVyIiwianRpIjoidGVzdDAwMSIsInN1YiI6InRlc3RUb3BpYyJ9LCJzIjoidGVzdERhdGEifQ.NLNmDo_Cq163gGEVFdQ_aFLNavnywlozw4XXhManHrE"))
	})
}

func (suite *JWTHS384Suite) TestHS384JWTParseTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewHS384JWT(suite.key)
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

func (suite *JWTHS384Suite) TestHS384JWTParseTokenMethodExpire() {
	suite.NotPanics(func() {
		j, err := NewHS384JWT(suite.key)
		suite.Equal(nil, err)

		standClaims := NewCommon(NewClaimsBuilder().Build())
		errParse := j.VerifyToken(
			"eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0ZXJDbGllbnQiLCJleHAiOjE2MzMzMjAzNDMsImlzcyI6InRlc3RlciIsImp0aSI6InRlc3QwMDEiLCJzIjoidGVzdERhdGEiLCJzdWIiOiJ0ZXN0VG9waWMifQ.ruXUXx-HMeuNXNuW2V2_SCj37PnnjLjinvq3gtxnGfkR8s_XC90M3_qEu3t4gfL5",
			standClaims)
		suite.Error(errParse, "token is expired")
	})
}

func (suite *JWTHS384Suite) TestHS384JWTParseTokenMethodExpireNoExpired() {
	suite.NotPanics(func() {
		j, err := NewHS384JWT(suite.key)
		suite.Equal(nil, err)

		standClaims := NewCommon(NewClaimsBuilder().Build())
		errParse := j.VerifyToken(
			"eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJJQ2xhaW1zIjp7ImF1ZCI6InRlc3RlckNsaWVudCIsImlzcyI6InRlc3RlciIsImp0aSI6InRlc3QwMDEiLCJzdWIiOiJ0ZXN0VG9waWMifSwicyI6InRlc3REYXRhIn0.PyiVJC6GI0XT7poIa_HJdTNzPG95i8laGAufTGxvkL6FXXnDtcwWBGP81Mrqk0lc",
			standClaims)
		suite.NoError(errParse)
	})
}

func (suite *JWTHS384Suite) TestHS384JWTRefreshTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewHS384JWT(suite.key)
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

func (suite *JWTHS384Suite) TestHS384JWTRefreshTokenMethodExpire() {
	suite.NotPanics(func() {
		j, err := NewHS384JWT(suite.key)
		suite.Equal(nil, err)
		standClaimsOutput := NewClaimsBuilder().Build()
		tk, errRefresh := j.RefreshToken(
			"eyJhbGciOiJIUzM4NCIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0ZXJDbGllbnQiLCJleHAiOjE2MzMzMjAzNDMsImlzcyI6InRlc3RlciIsImp0aSI6InRlc3QwMDEiLCJzIjoidGVzdERhdGEiLCJzdWIiOiJ0ZXN0VG9waWMifQ.ruXUXx-HMeuNXNuW2V2_SCj37PnnjLjinvq3gtxnGfkR8s_XC90M3_qEu3t4gfL5",
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

func (suite *JWTHS384Suite) TestHS384JWTRefreshTokenMethodparseRSRawError() {
	defer gostub.StubFunc(&parseSigned, nil, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		j, err := NewHS384JWT(suite.key)
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

func TestJWTHS384Suite(t *testing.T) {
	suite.Run(t, new(JWTHS384Suite))
}
