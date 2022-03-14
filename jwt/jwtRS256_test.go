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

type JWTRS256Suite struct {
	suite.Suite
	key    string
	option config.JWT
}

func (suite *JWTRS256Suite) SetupTest() {
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
}

func (suite *JWTRS256Suite) TestNewRS256JWT() {
	suite.NotPanics(func() {
		_, err := NewRS256JWT(suite.key)
		suite.Equal(nil, err)
	})
}

func (suite *JWTRS256Suite) TestNewRS256JWTFromOptions() {
	suite.NotPanics(func() {
		_, err := NewRS256JWTFromOptions(suite.option)
		suite.NoError(err)
	})
}

func (suite *JWTRS256Suite) TestNewRS256JWTFromOptionsAllNoOption() {
	suite.NotPanics(func() {
		_, err := NewRS256JWTFromOptions(config.JWT{
			RsaPrivateKeyPath: "",
			RsaPrivateKey:     "",
		})
		suite.Error(err, ErrNoKey)
	})
}

func (suite *JWTRS256Suite) TestNewRS256JWTFromOptionsAllNoOptionNoFIle() {
	suite.NotPanics(func() {
		_, err := NewRS256JWTFromOptions(config.JWT{
			RsaPrivateKeyPath: "testFile.txt",
			RsaPrivateKey:     "",
		})
		suite.Error(err)
	})
}

func (suite *JWTRS256Suite) TestNewRS256JWTparseRSAPrivateKeyFromPEMError() {
	defer gostub.StubFunc(&parseRSAPrivateKeyFromPEM, nil, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewRS256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWTRS256Suite) TestNewRS256JWTnewSignerError() {
	tSigner := testSigner{}
	defer gostub.StubFunc(&newSigner, tSigner, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewRS256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWTRS256Suite) TestRS256JWTGenerateTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewRS256JWT(suite.key)
		suite.Equal(nil, err)
		standClaims := NewClaimsBuilder().
			WithSubject("testTopic").
			WithIssuer("tester").
			WithID("test001").
			WithAudience([]string{"testerClient"}).
			ExpiresAfter(3 * time.Second).Build()
		//ExpiresAfter(87600 * time.Hour)
		common := NewCommon(standClaims, WithSecret("testData"))
		tk, errTk := j.GenerateToken(common)
		suite.T().Log(tk)
		suite.Equal(nil, errTk)
		suite.NotEmpty(tk)
	})
}

func (suite *JWTRS256Suite) TestRS256JWTValidateMethod() {
	suite.NotPanics(func() {
		j, err := NewRS256JWT(suite.key)
		suite.Equal(nil, err)

		suite.NoError(j.Validate(
			"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0ZXJDbGllbnQiLCJleHAiOjE2NDY5NjQ5NTYsImlzcyI6InRlc3RlciIsImp0aSI6InRlc3QwMDEiLCJzIjoidGVzdERhdGEiLCJzdWIiOiJ0ZXN0VG9waWMifQ.jLiy8FwQWMwEgOLaHVhZRunMuGPMX7zF9OjHsaI4zODgvu6RiVzqCpdh0rsWN33t7b327qWTzpU3r3cz_LN4wsEHVcqKd_wM-w4PtbMYzL1pOU7k8IyriEFHW8r3Zq9uynOEmBayWcsG9Dw_70xsmb2EAPT_0L8yR4quzqiTqDcqWC8h_uHycqBMw8CJhQVqiF0PqHWtoZTOGcomYnEWAIZHy7DVl-41jr7fc6jA9u1_g2EtvTP_DbwdgFawQ3ehTwkCCSGvxXdfPVvTvQJer43RRpSbqJ5ChZg4DiX83RD4F_JqQca8aaGS6kUcDdNbbqRRsLOowko3m3VO6bGPvg"))
	})
}

func (suite *JWTRS256Suite) TestRS256JWTValidateMethodparseSignedError() {
	suite.NotPanics(func() {
		defer gostub.StubFunc(&parseSigned, nil, errors.New("got error")).Reset()
		j, err := NewRS256JWT(suite.key)
		suite.Equal(nil, err)

		suite.Error(j.Validate(
			"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0ZXJDbGllbnQiLCJleHAiOjE2NDY5NjQ5NTYsImlzcyI6InRlc3RlciIsImp0aSI6InRlc3QwMDEiLCJzIjoidGVzdERhdGEiLCJzdWIiOiJ0ZXN0VG9waWMifQ.jLiy8FwQWMwEgOLaHVhZRunMuGPMX7zF9OjHsaI4zODgvu6RiVzqCpdh0rsWN33t7b327qWTzpU3r3cz_LN4wsEHVcqKd_wM-w4PtbMYzL1pOU7k8IyriEFHW8r3Zq9uynOEmBayWcsG9Dw_70xsmb2EAPT_0L8yR4quzqiTqDcqWC8h_uHycqBMw8CJhQVqiF0PqHWtoZTOGcomYnEWAIZHy7DVl-41jr7fc6jA9u1_g2EtvTP_DbwdgFawQ3ehTwkCCSGvxXdfPVvTvQJer43RRpSbqJ5ChZg4DiX83RD4F_JqQca8aaGS6kUcDdNbbqRRsLOowko3m3VO6bGPvg"))
	})
}

func (suite *JWTRS256Suite) TestRS256JWTParseTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewRS256JWT(suite.key)
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

func (suite *JWTRS256Suite) TestRS256JWTParseTokenMethodExpire() {
	suite.NotPanics(func() {
		j, err := NewRS256JWT(suite.key)
		suite.Equal(nil, err)

		standClaims := NewCommon(NewClaimsBuilder().Build())
		errParse := j.VerifyToken(
			"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0ZXJDbGllbnQiLCJleHAiOjE2MzMzMjAzODcsImlzcyI6InRlc3RlciIsImp0aSI6InRlc3QwMDEiLCJzIjoidGVzdERhdGEiLCJzdWIiOiJ0ZXN0VG9waWMifQ.Rc6Q6FA_Mo_WycfxQRVOktlGxoKU1BlmNMeo_Vrs1kfeSH_Ors63NhmoxAfMIGdEwOkO4kV56rLtRxRZ9ZUvyerB58QqbzeAhdRkXfwOIlc0GKSkwKuZK255akkPW0hFB7q1y839ZqukopaDRRZ7lUEmZl68t6KC63sNw2rdtHJwYkko6hQtVJsRlUJ0nREd59YHvKEMUsDcUFdTC24Kh5r3zCFG-vZC7Zl-YQsRKuZXE6wUNwOf3_RDu9rBKK1MpdE5zMv0_vlr4m7UJiMBSw22uildjh0nnD-a2ogfaCT2jwtcmi4AljStesxaUC8bjNM3MihLHv9dvC8lpkOWOw",
			standClaims)
		suite.Error(errParse, "token is expired")
	})
}

func (suite *JWTRS256Suite) TestRS256JWTParseTokenMethodExpireNoExpired() {
	suite.NotPanics(func() {
		j, err := NewRS256JWT(suite.key)
		suite.Equal(nil, err)

		standClaims := NewCommon(NewClaimsBuilder().Build())
		errParse := j.VerifyToken(
			"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0ZXJDbGllbnQiLCJpc3MiOiJ0ZXN0ZXIiLCJqdGkiOiJ0ZXN0MDAxIiwicyI6InRlc3REYXRhIiwic3ViIjoidGVzdFRvcGljIn0.Z8XUt7p41OJ92OafiahCkPaJ-fzm25-sxcJyD4MpOSgtRfa1QfhNllzeXEeIvwe4rEn2TZcc14p5tKHUelWz2UDBzXdoBKdO7XLE63wBN7zGMj_OmVddh1YrL_gguFndjPaAkawBFVcilvXKnbBWg58ZYFehtOxAcSxLgRhFt0E8qk6l351LBtlFNPJTm7RkL6SdYq9NO9i4otGxj_L2RabRrNJxTF07BCGTxuKQIH79udxTUgoB6BQxYHYQfUzD08UjGiWie93eP7HmuYqyyb4kESifYC6mD57aOSAQQbzSFxGNVfhTjIGEiLS6OfU_LwL2rTjRoYGl_4kRYM_MqA",
			standClaims)
		suite.NoError(errParse)
	})
}

func (suite *JWTRS256Suite) TestRS256JWTRefreshTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewRS256JWT(suite.key)
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

func (suite *JWTRS256Suite) TestRS256JWTRefreshTokenMethodExpire() {
	suite.NotPanics(func() {
		j, err := NewRS256JWT(suite.key)
		suite.Equal(nil, err)
		standClaimsOutput := NewClaimsBuilder().Build()
		tk, errRefresh := j.RefreshToken(
			"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0ZXJDbGllbnQiLCJleHAiOjE2NDY5NjU1NzEsImlzcyI6InRlc3RlciIsImp0aSI6InRlc3QwMDEiLCJzIjoidGVzdERhdGEiLCJzdWIiOiJ0ZXN0VG9waWMifQ.k0QMlXwXmOfAAbPOS8fY3tMla4Mb0bjB64-3BC2pQ4eo-paFs74sRiPSuIVX2NIyoVQD4L2BlJ5ly1uhR6WdLqdtaIAZ-A-bTSUHA26JU0t3LfAg2sfTEbHaff0TMUmhRNZ2EumCQJdNUiD8lv0g-ym2pRzaK0BFCqRj5OrEncUdc_9QtN8Pff9sWgv22vLBfK_teo5ZcIISVeBuje4huyVkLvdzeFKHAR6-jemRHXrnw6MvfdS65p4vnIpw1Hh21fM-kwCt_EbXaDFjWQOScNlJQsC_CIkD1q8MPJcxw3WsZcTNe7nG4_EFJjMkVDEVe87pO-7oMZqOZhlWNxxAWQ",
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

func (suite *JWTRS256Suite) TestRS256JWTRefreshTokenMethodparseRSRawError() {
	defer gostub.StubFunc(&parseSigned, nil, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		j, err := NewRS256JWT(suite.key)
		suite.Equal(nil, err)
		standClaimsOutput := NewClaimsBuilder().Build()
		_, errRefresh := j.RefreshToken(
			"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0ZXJDbGllbnQiLCJleHAiOjE2NDY5NjQ5NTYsImlzcyI6InRlc3RlciIsImp0aSI6InRlc3QwMDEiLCJzIjoidGVzdERhdGEiLCJzdWIiOiJ0ZXN0VG9waWMifQ.jLiy8FwQWMwEgOLaHVhZRunMuGPMX7zF9OjHsaI4zODgvu6RiVzqCpdh0rsWN33t7b327qWTzpU3r3cz_LN4wsEHVcqKd_wM-w4PtbMYzL1pOU7k8IyriEFHW8r3Zq9uynOEmBayWcsG9Dw_70xsmb2EAPT_0L8yR4quzqiTqDcqWC8h_uHycqBMw8CJhQVqiF0PqHWtoZTOGcomYnEWAIZHy7DVl-41jr7fc6jA9u1_g2EtvTP_DbwdgFawQ3ehTwkCCSGvxXdfPVvTvQJer43RRpSbqJ5ChZg4DiX83RD4F_JqQca8aaGS6kUcDdNbbqRRsLOowko3m3VO6bGPvg",
			standClaimsOutput,
			100*time.Millisecond,
		)
		suite.Error(errRefresh, "got error")
	})
}

func TestJWTRS256Suite(t *testing.T) {
	suite.Run(t, new(JWTRS256Suite))
}
