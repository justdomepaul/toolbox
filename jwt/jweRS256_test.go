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

type JWERS256Suite struct {
	suite.Suite
	key    string
	option config.JWT
	jwt    IJWT
}

func (suite *JWERS256Suite) SetupTest() {
	suite.key = `-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAvexSt0ZtW7qK9LBt73JgMISFBlFVxw54ukGOZvWCoSSEh12c
w67xy03nKlmADKac0ZPJ3K/TP2dMkLYi3fNojA//1wCkLhKdJBBZ+0i7qRtNXnSd
fHVWtFxXU7/zbSaRA+3YapiptGbL4OCRkgQaTnHKM+uQZt3yJeYeMSzF9hXwUs48
VJAdpt6AVZ6XDk4eC0D4Tyz9BXwU5oRFFn6qqdqWlLJ9E6o1ABr2atWYg26VmyPy
aMHBD4xQELT+Rmg3YsBonleGyA3j0yDedho+MUYrr8lJCKlv5nOI0tGhryw6XhRR
2ui/mq9xvVJ0jz6ozCR0cOldNHRLVF9aHDMbdwIDAQABAoIBAB756W/NA88fOMS5
9eRE8l1Xb97c6zGhMZ2nTZOLXXfs3dS6NvRPl05CcX6dxF3L2u3vvc/JuZmwvnMn
0b4Dkjyt61tk1mJRVOHp7NMoRLtLIa5TNNB0zuRx3yhguVJHJQXQCCkypxMuZPhT
iEqZcrTyqDkZpZ6xemomAyygEdWV7FAOsI1a7B/1qj+fm80orBs2Or8bgK+erclm
wNmelqp9KHP490MDgmQBeoYUItXLk2hx+VAiQCxQW+GVJlLLNpj1XgbAsyP0/zER
sddTnxbmaqlsOn4nxvyV130A5U4BnGVNeSgbbrdWE0QpJWzl8Ddnc4KxiitlbrT3
NYoF2tkCgYEA7bOp97HCuHy9Zc4WP2iz9b0EdpFitZW34S2hfHBdgD3vDMGXCmON
cyH75ZWYQh5M12Nm0zXmUPA0vp3lBh0hd1UP3f8tb2Hvyv8p2TyPqvtuQC6wAxuJ
ecO0EPWFVZmbGDItJpisBC1WCzo4AaW+NTJetrHHYKtzglX4Ts+OP5MCgYEAzIsY
eF3W3il82bz317x8pRZOLsuB6FQq3YZWuBEUP/l969koEF/QXxwttduqkyhQM22F
X1VzgyAGHCrdxEmGIteu/DBnJGNP6LwveO0O75isIEsqkSTA/zQ8ROch2cYglV9T
lDeQwKE6xzEpNj7+eRO1h3kHw8BomvPd2OdBOw0CgYBYW2V9viT8gNnCQwYAEgJ7
AQTssgQ4LWwJlvWlFPuclOkMG9XyNak5t9MztxS+1xaHJdrt/eYcBf4FMRoV2LQ8
8HCSe60+7u+8zHaY2qsoyodj8jbZIN5MVdPUTf9/HzcImnYwF6Yxc0y9pal16081
5QBR9ul+5JxuQVioqvxcYwKBgG++3hOEUMr2p3rdPhnio8YdNYFjNQmUUgbMSbwt
uH5q81xSOw0XC2OqpV5hMANNVuOBxgebS4wrhqsE0DtYX6vRYYvtdavvhcyEYvsR
p8NGCWNrLUo2ZioGg5axH1E2aL6yYZrr8G0MqGwCc51rNOM43Uex24gaKgvdhynk
zUJRAoGBAIz67SdKed8Xp+9x1D7aolIsgVLEcTkdJDkRigdbWN+0M89IGWo1o4Mz
y12fCTD1G9nbDrlMuYJ7TGedt65IRSqC5G3UMIa7r1pyLIn/K1eL06y5HMe9f5RO
XC+bv1Lzxf9y+eoDUaF/Q3U4h0BW4/tUF6BWVzM6/w4Dq7dqBCd2
-----END RSA PRIVATE KEY-----
`
	result := config.JWT{}
	suite.NoError(config.LoadFromEnv(&result))
	suite.option = result

	j, err := NewERS256JWT(suite.key)
	suite.NoError(err)
	suite.jwt = j
}

func (suite *JWERS256Suite) TestJWT() {
	suite.NotPanics(func() {
		_, err := NewERS256JWT(suite.key)
		suite.NoError(err)
	})
}

func (suite *JWERS256Suite) TestJWTFromOptions() {
	suite.NotPanics(func() {
		_, err := NewERS256JWTFromOptions(suite.option)
		suite.NoError(err)
	})
}

func (suite *JWERS256Suite) TestJWTFromOptionsAllNoOption() {
	suite.NotPanics(func() {
		_, err := NewERS256JWTFromOptions(config.JWT{
			RsaPrivateKeyPath: "",
			RsaPrivateKey:     "",
		})
		suite.Error(err, ErrNoKey)
	})
}

func (suite *JWERS256Suite) TestJWTFromOptionsAllNoOptionNoFIle() {
	suite.NotPanics(func() {
		_, err := NewERS256JWTFromOptions(config.JWT{
			RsaPrivateKeyPath: "testFile.txt",
			RsaPrivateKey:     "",
		})
		suite.Error(err)
	})
}

func (suite *JWERS256Suite) TestJWTParseRSAPrivateKeyFromPEMError() {
	defer gostub.StubFunc(&parseRSAPrivateKeyFromPEM, nil, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewERS256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWERS256Suite) TestJWTNewSignerError() {
	tSigner := testSigner{}
	defer gostub.StubFunc(&newSigner, tSigner, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewERS256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWERS256Suite) TestJWTNewEncrypterError() {
	tEncrypter := testEncrypter{}
	defer gostub.StubFunc(&newEncrypter, tEncrypter, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewERS256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWERS256Suite) TestJWTGenerateTokenMethod() {
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

func (suite *JWERS256Suite) TestJWTValidateMethod() {
	suite.NotPanics(func() {
		suite.NoError(suite.jwt.Validate(
			"eyJhbGciOiJSU0EtT0FFUC0yNTYiLCJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2R0NNIiwidHlwIjoiSldUIn0.buc8jTfAMFU2TcC4rUSeU6QxwfMHMwsjmixztdIKjGyTvp7c8jQLlO8FyGF_KnVw1CqzI79jPJEaTQOGPBZjsbe_bkNjKiAc-I8Hc8yxt2u401tClvfYSDFFTJMziSv5CyShfJ9kT2q779p7PjXUgp0Bn6yARX5lhI2sgbFYRgp6c5WVo7k7woGE9GeBbtJK9q0QJo5E94xVyVzcfTkNEsgvG1XQz2e4HyMp6ug_BYnYPrQkjdLrqqpRo4ApQ6wIIcm5O-Fj0fulGpVyEX4sdL1TMJsW2vO5Jb7upjOjMyYFVUN5SoyOqzH7gx7GZLFG30vh9uB8rx6bViSzJVct7g.S5kq0Z450P4JA9SZ.JFPnRvQ6sroSfPECKgkbglLxoz7COsl8GR039vFOKCSwifuPKeoxcunD89cNU4PbBKmHLI3DCDfeHrQn1n7xzQnSPAewuz6BuDKDRd06GCGeKrUu2s3FA57DGPhswgoEbFCxuh62osqOPCz9HoYph9RuMvtEqPJzsmP82Va8lQb5XdiiVwHRNCOPFMzzEgPDhFryRLKYhKm-LdLAHwZH-t4t6DNU7N1D2iB1lWdVTExPdvdDwZpN05zzJZw2Cotx-XKeHiZ7FqzoocikiXCHqPI3jlHIhZFJoIew46PLH98j5i3F7QihEyVVKdSuHc3ZVeoNT6lexjkbrl-ofxp08dWfeXla8XhybLG70RhDdc-VmrGW0nOX16QTKMCHssr05mseYzYpsDXaJQxCBPuKYeAmCZJrOKeX4fdqICi3vLtgC-LgmBKkFutKqAIN3wQoGjJmN0t2fSNW1PU8dfYx39WPYUeZxYG9JJQ98LPB5u6AENdTcFd6KfcSMft1OdoIlTMATG3Gnxe9H7Af3L98F796oOlPi1G9DNWQKyXpbfX6VytRGuQaiMITPbyrw0zLnyYI3gCgDwndbX1nWV-E-8E4cp_vpk-hTprn39kN0kNfBxaZJWQ3zjIB9dwhWWQ3DTkgwdJe237OaNJcBNTfGeULW2fAvObO4QRwKmdBOHWSdYoUSss.4Ur11EeA5cv7tEHTI4OP4g"))
	})
}

func (suite *JWERS256Suite) TestJWTValidateMethodParseSignedAndEncryptedError() {
	suite.NotPanics(func() {
		defer gostub.StubFunc(&parseSignedAndEncrypted, nil, errors.New("got error")).Reset()

		suite.Error(suite.jwt.Validate(""))
	})
}

func (suite *JWERS256Suite) TestJWTVerifyTokenMethod() {
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

func (suite *JWERS256Suite) TestJWTVerifyTokenMethodExpire() {
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

func (suite *JWERS256Suite) TestJWTVerifyTokenMethodExpireNoExpired() {
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

func (suite *JWERS256Suite) TestJWTRefreshTokenMethod() {
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

func (suite *JWERS256Suite) TestJWTRefreshTokenMethodExpire() {
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

func (suite *JWERS256Suite) TestJWTRefreshTokenMethodExpireNoIJWTExpireClaim() {
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

func (suite *JWERS256Suite) TestJWTRefreshTokenMethodParseRawError() {
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

func TestJWERS256Suite(t *testing.T) {
	suite.Run(t, new(JWERS256Suite))
}
