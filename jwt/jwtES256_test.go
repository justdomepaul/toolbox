package jwt

import (
	"github.com/cockroachdb/errors"
	"github.com/justdomepaul/toolbox/config"
	"github.com/prashantv/gostub"
	"github.com/square/go-jose/v3/jwt"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type JWTES256Suite struct {
	suite.Suite
	key    string
	option config.JWT
	jwt    IJWT
}

func (suite *JWTES256Suite) SetupTest() {
	suite.key = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIANpsaiXN0UdToO728+5kwGNPA+RsbEPwb3MGiAlBhSXoAoGCCqGSM49
AwEHoUQDQgAER5WNvPs/SMICGESgDbN7IYl0CvPSkUhAaUtF/LAQEINqte/HLMks
hRsKJ2MTCe1upn5vhgBuGl5CL4ea4DqNhA==
-----END EC PRIVATE KEY-----
`
	result := config.JWT{}
	suite.NoError(config.LoadFromEnv(&result))
	suite.option = result

	j, err := NewES256JWT(suite.key)
	suite.NoError(err)
	suite.jwt = j
}

func (suite *JWTES256Suite) TestNewJWT() {
	suite.NotPanics(func() {
		_, err := NewES256JWT(suite.key)
		suite.NoError(err)
	})
}

func (suite *JWTES256Suite) TestJWTFromOptions() {
	suite.NotPanics(func() {
		_, err := NewES256JWTFromOptions(suite.option)
		suite.NoError(err)
	})
}

func (suite *JWTES256Suite) TestJWTFromOptionsAllNoOption() {
	suite.NotPanics(func() {
		_, err := NewES256JWTFromOptions(config.JWT{
			EcdsaPrivateKeyPath: "",
			EcdsaPrivateKey:     "",
		})
		suite.Error(err)
	})
}

func (suite *JWTES256Suite) TestJWTFromOptionsAllNoOptionNoFIle() {
	suite.NotPanics(func() {
		_, err := NewES256JWTFromOptions(config.JWT{
			EcdsaPrivateKeyPath: "testFile.txt",
			EcdsaPrivateKey:     "",
		})
		suite.Error(err)
	})
}

func (suite *JWTES256Suite) TestJWTParseECPrivateKeyFromPEMError() {
	defer gostub.StubFunc(&parseECPrivateKeyFromPEM, nil, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewES256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWTES256Suite) TestJWTNewSignerError() {
	tSigner := testSigner{}
	defer gostub.StubFunc(&newSigner, tSigner, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewES256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWTES256Suite) TestJWTGenerateTokenMethod() {
	suite.NotPanics(func() {
		standClaims := NewClaimsBuilder().
			WithSubject("testTopic").
			WithIssuer("tester").
			WithID("test001").
			WithAudience([]string{"testerClient"}).
			ExpiresAfter(5 * time.Second).Build()
		common := NewCommon(standClaims, WithSecret("testData"))
		tk, errTk := suite.jwt.GenerateToken(common)
		suite.T().Log(tk)
		suite.NoError(errTk)
		suite.NotEmpty(tk)
	})
}

func (suite *JWTES256Suite) TestJWTValidateMethod() {
	suite.NotPanics(func() {
		suite.NoError(suite.jwt.Validate(
			"eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0ZXJDbGllbnQiLCJleHAiOjE2NDc0MDQ0NzQsImlzcyI6InRlc3RlciIsImp0aSI6InRlc3QwMDEiLCJzIjoidGVzdERhdGEiLCJzdWIiOiJ0ZXN0VG9waWMifQ.yW0MtToyyaRfSRVXMWLWoG_2hbsMh_r2HfKCUF4oXkXwvj0H_kZ7PCoGd03AEx_kjsuJqo8uptOryHUFqN6Mww"))
	})
}

func (suite *JWTES256Suite) TestJWTValidateMethodParseSignedError() {
	suite.NotPanics(func() {
		defer gostub.StubFunc(&parseSigned, nil, errors.New("got error")).Reset()
		suite.Error(suite.jwt.Validate(""))
	})
}

func (suite *JWTES256Suite) TestJWTVerifyTokenMethod() {
	suite.NotPanics(func() {
		standClaims := NewClaimsBuilder().
			WithSubject("testTopic").
			WithIssuer("tester").
			WithID("test001").
			WithAudience([]string{"testerClient"}).
			ExpiresAfter(3 * time.Second).Build()
		common := NewCommon(standClaims, WithSecret("testData"))
		tk, errTk := suite.jwt.GenerateToken(common)
		suite.NoError(errTk)
		suite.T().Log(tk)
		standCommonOutput := NewCommon(NewClaimsBuilder().Build())
		errParse := suite.jwt.VerifyToken(
			tk,
			standCommonOutput)
		suite.NoError(errParse)
		suite.Equal("testTopic", standCommonOutput.Subject)
		suite.Equal("tester", standCommonOutput.Issuer)
		suite.Equal("test001", standCommonOutput.ID)
		suite.Equal(jwt.Audience{"testerClient"}, standCommonOutput.Audience)
	})
}

func (suite *JWTES256Suite) TestJWTVerifyTokenMethodExpire() {
	suite.NotPanics(func() {
		standClaims := NewClaimsBuilder().
			WithSubject("testTopic").
			WithIssuer("tester").
			WithID("test001").
			WithAudience([]string{"testerClient"}).
			ExpiresAfter(-3 * time.Second).Build()
		common := NewCommon(standClaims, WithSecret("testData"))
		tk, errTk := suite.jwt.GenerateToken(common)
		suite.NoError(errTk)
		suite.T().Log(tk)
		standCommonOutput := NewCommon(NewClaimsBuilder().Build())
		errParse := suite.jwt.VerifyToken(
			tk,
			standCommonOutput)
		suite.ErrorIs(errParse, ErrTokenExpired)
	})
}

func (suite *JWTES256Suite) TestJWTVerifyTokenMethodExpireNoExpired() {
	suite.NotPanics(func() {
		standClaims := NewClaimsBuilder().
			WithSubject("testTopic").
			WithIssuer("tester").
			WithID("test001").
			WithAudience([]string{"testerClient"}).Build()
		common := NewCommon(standClaims, WithSecret("testData"))
		tk, errTk := suite.jwt.GenerateToken(common)
		suite.NoError(errTk)
		suite.T().Log(tk)
		standCommonOutput := NewClaimsBuilder().Build()
		errParse := suite.jwt.VerifyToken(
			tk,
			standCommonOutput)
		suite.NoError(errParse)
		suite.Equal("testTopic", standCommonOutput.Subject)
		suite.Equal("tester", standCommonOutput.Issuer)
		suite.Equal("test001", standCommonOutput.ID)
		suite.Equal(jwt.Audience{"testerClient"}, standCommonOutput.Audience)
	})
}

func (suite *JWTES256Suite) TestJWTRefreshTokenMethod() {
	suite.NotPanics(func() {
		standClaims := NewClaimsBuilder().
			WithSubject("testTopic").
			WithIssuer("tester").
			WithID("test001").
			WithAudience([]string{"testerClient"}).
			ExpiresAfter(-3 * time.Second).Build()
		common := NewCommon(standClaims, WithSecret("testData"))
		tk, errTk := suite.jwt.GenerateToken(common)
		suite.NoError(errTk)
		suite.T().Log(tk)
		standCommonOutput := NewCommon(NewClaimsBuilder().Build())
		newTk, errRefresh := suite.jwt.RefreshToken(
			tk,
			standCommonOutput,
			100*time.Millisecond)
		suite.NoError(errRefresh)
		suite.NotEqual(tk, newTk)
	})
}

func (suite *JWTES256Suite) TestJWTRefreshTokenMethodExpire() {
	suite.NotPanics(func() {
		standClaims := NewClaimsBuilder().
			WithSubject("testTopic").
			WithIssuer("tester").
			WithID("test001").
			WithAudience([]string{"testerClient"}).
			ExpiresAfter(-3 * time.Second).Build()
		common := NewCommon(standClaims, WithSecret("testData"))
		tk, errTk := suite.jwt.GenerateToken(common)
		suite.NoError(errTk)
		suite.T().Log(tk)
		standClaimsOutput := NewCommon(NewClaimsBuilder().Build())
		tkRefresh, errRefresh := suite.jwt.RefreshToken(
			tk,
			standClaimsOutput,
			100*time.Millisecond,
		)
		suite.NoError(errRefresh)
		suite.NotEmpty(tk)
		suite.NotEqual(tk, tkRefresh)
	})
}

func (suite *JWTES256Suite) TestJWTRefreshTokenMethodExpireNoIJWTExpireClaim() {
	suite.NotPanics(func() {
		standClaims := NewClaimsBuilder().
			WithSubject("testTopic").
			WithIssuer("tester").
			WithID("test001").
			WithAudience([]string{"testerClient"}).
			ExpiresAfter(-3 * time.Second).Build()
		common := NewMockClaim(standClaims)
		tk, errTk := suite.jwt.GenerateToken(common)
		suite.NoError(errTk)
		suite.T().Log(tk)
		standClaimsOutput := NewMockClaim(NewClaimsBuilder().Build())
		tkRefresh, errRefresh := suite.jwt.RefreshToken(
			tk,
			standClaimsOutput,
			100*time.Millisecond,
		)
		suite.NoError(errRefresh)
		suite.NotEmpty(tk)
		suite.Equal(tk, tkRefresh)
	})
}

func (suite *JWTES256Suite) TestJWTRefreshTokenMethodParseRawError() {
	defer gostub.StubFunc(&parseSigned, nil, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		standClaimsOutput := NewClaimsBuilder().Build()
		_, errRefresh := suite.jwt.RefreshToken(
			"",
			standClaimsOutput,
			100*time.Millisecond,
		)
		suite.Error(errRefresh, "got error")
	})
}

func TestJWTES256Suite(t *testing.T) {
	suite.Run(t, new(JWTES256Suite))
}
