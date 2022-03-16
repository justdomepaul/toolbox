package jwt

import (
	"github.com/cockroachdb/errors"
	"github.com/justdomepaul/toolbox/config"
	"github.com/justdomepaul/toolbox/key"
	"github.com/prashantv/gostub"
	"github.com/square/go-jose/v3/jwt"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type JWEHS384Suite struct {
	suite.Suite
	key    string
	option config.JWT
	jwt    IJWT
}

func (suite *JWEHS384Suite) SetupTest() {
	suite.key = `GhZ6ZhYNJijb-VNQ0H9EiIzF4JiQbJqryDNTuALnPTRIFaLid2091W0wuqAD9kcm`
	result := config.JWT{}
	suite.NoError(config.LoadFromEnv(&result))
	suite.option = result

	j, err := NewEHS384JWT(key.ToBinaryRunes(suite.key))
	suite.NoError(err)
	suite.jwt = j
}

func (suite *JWEHS384Suite) TestJWT() {
	suite.NotPanics(func() {
		_, err := NewEHS384JWT(suite.key)
		suite.NoError(err)
	})
}

func (suite *JWEHS384Suite) TestJWTFromOptions() {
	suite.NotPanics(func() {
		_, err := NewEHS384JWTFromOptions(suite.option)
		suite.NoError(err)
	})
}

func (suite *JWEHS384Suite) TestJWTFromOptionsAllNoOption() {
	suite.NotPanics(func() {
		_, err := NewEHS384JWTFromOptions(config.JWT{
			HmacSecretKeyPath: "",
			HmacSecretKey:     "",
		})
		suite.Error(err, ErrNoKey)
	})
}

func (suite *JWEHS384Suite) TestJWTFromOptionsAllNoOptionNoFIle() {
	suite.NotPanics(func() {
		_, err := NewEHS384JWTFromOptions(config.JWT{
			HmacSecretKeyPath: "testFile.txt",
			HmacSecretKey:     "",
		})
		suite.Error(err)
	})
}

func (suite *JWEHS384Suite) TestJWTNewSignerError() {
	tSigner := testSigner{}
	defer gostub.StubFunc(&newSigner, tSigner, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewEHS384JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWEHS384Suite) TestJWTnewEncrypterError() {
	tEncrypter := testEncrypter{}
	defer gostub.StubFunc(&newEncrypter, tEncrypter, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewEHS384JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWEHS384Suite) TestJWTGenerateTokenMethod() {
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

func (suite *JWEHS384Suite) TestJWTValidateMethod() {
	suite.NotPanics(func() {
		suite.NoError(suite.jwt.Validate(
			"eyJhbGciOiJQQkVTMi1IUzM4NCtBMTkyS1ciLCJjdHkiOiJKV1QiLCJlbmMiOiJBMTkyQ0JDLUhTMzg0IiwicDJjIjoxMDAwMDAsInAycyI6IlMwY0ZPMU1yMHVTRkJzUm1MM2ZoVmciLCJ0eXAiOiJKV1QifQ.nO0L-I9v5sXr158ltqwoNaeXoWdCIjZ2G0rSLfjLXMCThAPSFUSnDiVrUxtF92-Owmj_-4i2JT4.aPHgG1bKpDqVu1N_orI0Lg.yqGNpYmsTmxf5awYiiL85kG1KF2JkMrMFzg_WXtQ8j5YpPX8dGORVZ5wOCWGu-X7cHTGAYoEgIYmRE7ZjBC4QHnkeraB74Zlgctz7NyVF9qQRpihHmasXuacxQb6xm-P5O-BvUhXxO2gE6R3xrc1IPq1hO1PyWbT0eq2X2wQfDKE2nR8gA29XzX3K26srSXSZqVo6xu9nzyr-wr9-1XiIs2sqNKCZtXxAJPHPbECgYAaIA74q4Iz-2mTkvwP2hHRcxUNcqFeS31KmD74yO7spOgsgMfxd1kTORqoAD3hIAauSg7N8cSEUz7uDXl93dCAE8kPQK_0Dmf4KRREPl-paA.7RbTwjNxJRQo8SkxfEygFnlcTbiIhlqX"))
	})
}

func (suite *JWEHS384Suite) TestJWTValidateMethodparseSignedAndEncryptedError() {
	suite.NotPanics(func() {
		defer gostub.StubFunc(&parseSignedAndEncrypted, nil, errors.New("got error")).Reset()

		suite.Error(suite.jwt.Validate(""))
	})
}

func (suite *JWEHS384Suite) TestJWTVerifyTokenMethod() {
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

func (suite *JWEHS384Suite) TestJWTVerifyTokenMethodExpire() {
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

func (suite *JWEHS384Suite) TestJWTVerifyTokenMethodExpireNoExpired() {
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

func (suite *JWEHS384Suite) TestJWTRefreshTokenMethod() {
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
		newTk, errRefresh := suite.jwt.RefreshToken(
			tk,
			standCommonOutput,
			100*time.Millisecond)
		suite.NoError(errRefresh)
		suite.Equal(tk, newTk)
	})
}

func (suite *JWEHS384Suite) TestJWTRefreshTokenMethodExpire() {
	suite.NotPanics(func() {
		standClaims := NewClaimsBuilder().
			WithSubject("testTopic").
			WithIssuer("tester").
			WithID("test001").
			WithAudience([]string{"testerClient"}).
			ExpiresAfter(-3 * time.Second).Build()
		common := NewCommon(standClaims, WithSecret("testData"))
		tk, errTk := suite.jwt.GenerateToken(common)
		suite.T().Log(tk)
		suite.NoError(errTk)
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

func (suite *JWEHS384Suite) TestJWTRefreshTokenMethodExpireNoIJWTExpireClaim() {
	suite.NotPanics(func() {
		standClaims := NewClaimsBuilder().
			WithSubject("testTopic").
			WithIssuer("tester").
			WithID("test001").
			WithAudience([]string{"testerClient"}).
			ExpiresAfter(-3 * time.Second).Build()
		common := NewMockClaim(standClaims)
		tk, errTk := suite.jwt.GenerateToken(common)
		suite.T().Log(tk)
		suite.NoError(errTk)
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

func (suite *JWEHS384Suite) TestJWTRefreshTokenMethodparseRawError() {
	defer gostub.StubFunc(&parseSignedAndEncrypted, nil, errors.New("got error")).Reset()
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

func TestJWEHS384Suite(t *testing.T) {
	suite.Run(t, new(JWEHS384Suite))
}
