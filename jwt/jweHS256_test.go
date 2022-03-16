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

type JWEHS256Suite struct {
	suite.Suite
	key    string
	option config.JWT
	jwt    IJWT
}

func (suite *JWEHS256Suite) SetupTest() {
	suite.key = `b583ed184e2018b3d89a4fa8832d0a1f`
	result := config.JWT{}
	suite.NoError(config.LoadFromEnv(&result))
	suite.option = result

	j, err := NewEHS256JWT(key.ToBinaryRunes(suite.key))
	suite.NoError(err)
	suite.jwt = j
}

func (suite *JWEHS256Suite) TestJWT() {
	suite.NotPanics(func() {
		_, err := NewEHS256JWT(suite.key)
		suite.NoError(err)
	})
}

func (suite *JWEHS256Suite) TestJWTFromOptions() {
	suite.NotPanics(func() {
		_, err := NewEHS256JWTFromOptions(suite.option)
		suite.NoError(err)
	})
}

func (suite *JWEHS256Suite) TestJWTFromOptionsAllNoOption() {
	suite.NotPanics(func() {
		_, err := NewEHS256JWTFromOptions(config.JWT{
			HmacSecretKeyPath: "",
			HmacSecretKey:     "",
		})
		suite.Error(err, ErrNoKey)
	})
}

func (suite *JWEHS256Suite) TestJWTFromOptionsAllNoOptionNoFIle() {
	suite.NotPanics(func() {
		_, err := NewEHS256JWTFromOptions(config.JWT{
			HmacSecretKeyPath: "testFile.txt",
			HmacSecretKey:     "",
		})
		suite.Error(err)
	})
}

func (suite *JWEHS256Suite) TestJWTNewSignerError() {
	tSigner := testSigner{}
	defer gostub.StubFunc(&newSigner, tSigner, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewEHS256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWEHS256Suite) TestJWTNewEncrypterError() {
	tEncrypter := testEncrypter{}
	defer gostub.StubFunc(&newEncrypter, tEncrypter, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewEHS256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWEHS256Suite) TestJWTGenerateTokenMethod() {
	suite.NotPanics(func() {
		standClaims := NewClaimsBuilder().
			WithSubject("testTopic").
			WithIssuer("tester").
			WithID("test001").
			WithAudience([]string{"testerClient"}).
			ExpiresAfter(3 * time.Second).Build()
		common := NewCommon(standClaims, WithSecret("testData"))
		tk, errTk := suite.jwt.GenerateToken(common)
		suite.T().Log(tk)
		suite.NoError(errTk)
		suite.NotEmpty(tk)
	})
}

func (suite *JWEHS256Suite) TestJWTValidateMethod() {
	suite.NotPanics(func() {
		suite.NoError(suite.jwt.Validate(
			"eyJhbGciOiJQQkVTMi1IUzI1NitBMTI4S1ciLCJjdHkiOiJKV1QiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2IiwicDJjIjoxMDAwMDAsInAycyI6Im91T25zQmJOSFlwRDZfVHBzcmRYc2ciLCJ0eXAiOiJKV1QifQ.IqAAjhKuYBH4R5WVTgm4-8r-JRCoEez1VXL315LmedhlRIG8rA8KRw.z8FcZmitQlr7Dwv8jDpJVA.kpMyD3xCQ41MmmegHZ7V6sLEPUyyhX2P3MebaSTvmS7dq7EcEZ670kfg64fAlnRh1gm_vZ2wVjkHCRoEDhKcrX05U2I7IpkekXsHMFxrCeC_CJcIUiVPBvbgX7zm_vJVpw_ihBJw7UK3UNp2vUrYjRy--7TZdurJlx1JsQ9QePjHWp3B9_Q5IJdj47AArIqgTExwWO27Xs-MeSCP05__0CSsWtW20E3NsOSPUWfV-rmpyk7weDM2L3Bf7O6F4ALuew6ybON_ahE0cv5XcDRgdw36S3hwO2JTaL_DQUbXaTw.96OEmUK272U37oeOxX91qg"))
	})
}

func (suite *JWEHS256Suite) TestJWTValidateMethodParseSignedAndEncryptedError() {
	suite.NotPanics(func() {
		defer gostub.StubFunc(&parseSignedAndEncrypted, nil, errors.New("got error")).Reset()

		suite.Error(suite.jwt.Validate(""))
	})
}

func (suite *JWEHS256Suite) TestJWTVerifyTokenMethod() {
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

func (suite *JWEHS256Suite) TestJWTVerifyTokenMethodExpire() {
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

func (suite *JWEHS256Suite) TestJWTVerifyTokenMethodExpireNoExpiredTime() {
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

func (suite *JWEHS256Suite) TestJWTRefreshTokenMethod() {
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

func (suite *JWEHS256Suite) TestJWTRefreshTokenMethodExpire() {
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

func (suite *JWEHS256Suite) TestJWTRefreshTokenMethodExpireNoIJWTExpireClaim() {
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

func (suite *JWEHS256Suite) TestJWTRefreshTokenMethodParseRawError() {
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

func TestJWEHS256Suite(t *testing.T) {
	suite.Run(t, new(JWEHS256Suite))
}
