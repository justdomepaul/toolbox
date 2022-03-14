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
}

func (suite *JWERS256Suite) TestNewERS256JWT() {
	suite.NotPanics(func() {
		_, err := NewERS256JWT(suite.key)
		suite.Equal(nil, err)
	})
}

func (suite *JWERS256Suite) TestNewERS256JWTFromOptions() {
	suite.NotPanics(func() {
		_, err := NewERS256JWTFromOptions(suite.option)
		suite.NoError(err)
	})
}

func (suite *JWERS256Suite) TestNewERS256JWTFromOptionsAllNoOption() {
	suite.NotPanics(func() {
		_, err := NewERS256JWTFromOptions(config.JWT{
			RsaPrivateKeyPath: "",
			RsaPrivateKey:     "",
		})
		suite.Error(err, ErrNoKey)
	})
}

func (suite *JWERS256Suite) TestNewERS256JWTFromOptionsAllNoOptionNoFIle() {
	suite.NotPanics(func() {
		_, err := NewERS256JWTFromOptions(config.JWT{
			RsaPrivateKeyPath: "testFile.txt",
			RsaPrivateKey:     "",
		})
		suite.Error(err)
	})
}

func (suite *JWERS256Suite) TestNewERS256JWTparseRSAPrivateKeyFromPEMError() {
	defer gostub.StubFunc(&parseRSAPrivateKeyFromPEM, nil, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewERS256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWERS256Suite) TestNewERS256JWTnewSignerError() {
	tSigner := testSigner{}
	defer gostub.StubFunc(&newSigner, tSigner, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewERS256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWERS256Suite) TestNewERS256JWTnewEncrypterError() {
	tEncrypter := testEncrypter{}
	defer gostub.StubFunc(&newEncrypter, tEncrypter, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewERS256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWERS256Suite) TestERS256JWTGenerateTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewERS256JWT(suite.key)
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

func (suite *JWERS256Suite) TestERS256JWTValidateMethod() {
	suite.NotPanics(func() {
		j, err := NewERS256JWT(suite.key)
		suite.Equal(nil, err)

		suite.NoError(j.Validate(
			"eyJhbGciOiJSU0EtT0FFUC0yNTYiLCJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2R0NNIiwidHlwIjoiSldUIn0.W81uCiTrZ1j-w311Q512RshQFpAIlXxlGnhlQVITwewfY8DoF0yJJgXQ8ifCvCu7n-IvD2cB7K2fmfiiCcJTsbKipHJg0DT9MxawV_JhaBRrzkKtzR842G4s8Z9EadefMcUFdaZ_klQ7I7P1ieMjDfd-HYKzleOip3BRRe0Y33AP803X5K584pTTUW9_vR2Vc305R1IF4qx9OpJiSZmNEZ-wP23WVhsegkBruDgPf0s4yJsVV1rdUZ9JFR9XNjq_6VvZ4klo1ODHqIZ1-4V3Gs1y8ShQKx2cb29gjoqtV01-7ik-CFWYmL83c6oTv_mjHPyCboVfrIAtIra1sJqOPA.tqaK3VNvv5_IZuwF.sdnCs6ovZICVQ50ypywZgEYeS7pgPHxdYoIBIg0Wm10LEpPCzcYT86ToufVzR4_IHGeHGZa7fiLGUJy4IzOvklhSjNGmqwWnZf5jfXQBwgizrTT0ZsdqnEHXEQWN4ATOoQ4kAnJRoh9gfJoOObPAilCWgvfZjE9Kq3cYNjbY7SJPhujmSLT5jwWiHx4-TbTGVU6CRmGUe5_EnfDvGtsxf0E1Pndevcu1aW5eOn6qYTFHX1H2NImuMBUc2sradY6KNPm_udd_qEs3jCCsGoDvtlLv3yEVhpkx1bHWiQdHx-2ZKnd4o8lI2jzjwojbae3VNqXaCaHjjcI2crpeYGKLlQ-xJpIVMRo3CCJZunoO_0dL6qLEIyeB_W7QfnWQBcddV1K8r5bbCU0AJbz6AqMNci51DO2_zCQEXVSMJCAAhdp8jo39leb8kfcgOkllZxi1P3xUexsy5r8ZwmuoB6s9Ou1-k1-UPUUen2XJhCtfB7-Hrz5BoxyrWsbtfLF22G3QAw-HwYltyrMbqpp8monT10usPZ2nJpp9F3H5AzdyGKphcJLB1uUll-c1Iu8RIGfcrch_ZuWAYDvn9QdpXFJNE_pgN6fUGYg7YXKBmYyNUlOZRsNrMMrg3x5L5q8WFGCecjZWjpXjjIeB-bnPISpH1zx22_stSsbjuIRjQA1F8-v75wvReO4.ns0TF_uNBL0HKBmvM2tEvQ"))
	})
}

func (suite *JWERS256Suite) TestERS256JWTValidateMethodparseSignedAndEncryptedError() {
	suite.NotPanics(func() {
		defer gostub.StubFunc(&parseSignedAndEncrypted, nil, errors.New("got error")).Reset()
		j, err := NewERS256JWT(suite.key)
		suite.Equal(nil, err)

		suite.Error(j.Validate(
			"eyJhbGciOiJSU0EtT0FFUC0yNTYiLCJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2R0NNIiwidHlwIjoiSldUIn0.W81uCiTrZ1j-w311Q512RshQFpAIlXxlGnhlQVITwewfY8DoF0yJJgXQ8ifCvCu7n-IvD2cB7K2fmfiiCcJTsbKipHJg0DT9MxawV_JhaBRrzkKtzR842G4s8Z9EadefMcUFdaZ_klQ7I7P1ieMjDfd-HYKzleOip3BRRe0Y33AP803X5K584pTTUW9_vR2Vc305R1IF4qx9OpJiSZmNEZ-wP23WVhsegkBruDgPf0s4yJsVV1rdUZ9JFR9XNjq_6VvZ4klo1ODHqIZ1-4V3Gs1y8ShQKx2cb29gjoqtV01-7ik-CFWYmL83c6oTv_mjHPyCboVfrIAtIra1sJqOPA.tqaK3VNvv5_IZuwF.sdnCs6ovZICVQ50ypywZgEYeS7pgPHxdYoIBIg0Wm10LEpPCzcYT86ToufVzR4_IHGeHGZa7fiLGUJy4IzOvklhSjNGmqwWnZf5jfXQBwgizrTT0ZsdqnEHXEQWN4ATOoQ4kAnJRoh9gfJoOObPAilCWgvfZjE9Kq3cYNjbY7SJPhujmSLT5jwWiHx4-TbTGVU6CRmGUe5_EnfDvGtsxf0E1Pndevcu1aW5eOn6qYTFHX1H2NImuMBUc2sradY6KNPm_udd_qEs3jCCsGoDvtlLv3yEVhpkx1bHWiQdHx-2ZKnd4o8lI2jzjwojbae3VNqXaCaHjjcI2crpeYGKLlQ-xJpIVMRo3CCJZunoO_0dL6qLEIyeB_W7QfnWQBcddV1K8r5bbCU0AJbz6AqMNci51DO2_zCQEXVSMJCAAhdp8jo39leb8kfcgOkllZxi1P3xUexsy5r8ZwmuoB6s9Ou1-k1-UPUUen2XJhCtfB7-Hrz5BoxyrWsbtfLF22G3QAw-HwYltyrMbqpp8monT10usPZ2nJpp9F3H5AzdyGKphcJLB1uUll-c1Iu8RIGfcrch_ZuWAYDvn9QdpXFJNE_pgN6fUGYg7YXKBmYyNUlOZRsNrMMrg3x5L5q8WFGCecjZWjpXjjIeB-bnPISpH1zx22_stSsbjuIRjQA1F8-v75wvReO4.ns0TF_uNBL0HKBmvM2tEvQ"))
	})
}

func (suite *JWERS256Suite) TestERS256JWTParseTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewERS256JWT(suite.key)
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

func (suite *JWERS256Suite) TestERS256JWTParseTokenMethodExpire() {
	suite.NotPanics(func() {
		j, err := NewERS256JWT(suite.key)
		suite.Equal(nil, err)

		standCommonOutput := NewCommon(NewClaimsBuilder().Build())
		errParse := j.VerifyToken(
			"eyJhbGciOiJSU0EtT0FFUC0yNTYiLCJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2R0NNIiwidHlwIjoiSldUIn0.SyZFN3mJUA5LT10gmLd4APZR07Yff4UBPLM9TrRnwwGCVBWP5v2W_Kt32eoV0ogcPwXEKNqs_x4ymxkzAzb8PUaaUDynR0glUiN41aJi-pDT-eT_tGorUpv7jFXHaUJvIcrXJUVuEs3udgFqXOvtwVvC8lKPmLfR3C_aLBG1lOIoO67dejavsh1Dr4tJuaeduCctRaq-kndYLaTokGkSgQ6furgPCiHNc_HXdgVUL_SZQ9oINAwBvio3B2D6fwltIjpEgu_5REYAxsYX97IrTVvmv5TLXV360Qg8skWyn_tkJSDlgArAG9HuDe2CQ1QOVsRuQu6P8f019AyEaG_Peg.OFdrCq2mbVLQpCQF.ifBTlPiz1lorEuw0C61Wsq5ajB3fsoS1E5-qeGTw3N3UztYnEsgfEaA0gkszMChO1kOiER0y2cyKnZNjMRyjsK8QC9itxCjQSx6VQY3WVqUpLFGUy692paeRHTqYLUkHyr5xezTFyPxv0i_uU9xw7u6T5XtbsvtEmLyzgNgm7cMjyXN7zeSyN5lWwoUzkrdNpmgSzzR2Q7fpD_8F8osOEysIOSbWAEsokOldFPM_xoHq6PZb_nN9ugbzqZbl8NwImK030BjrXVVZ3P8qY-cBID6QBeq77hGjUfdh5ctWWd7NVgSfmKs8Gm1tNRuaE1gek3fTyJcJZHxap6t_1fsPX7sc0l93GwX-9mw4lubkg-RlymZvzoG5EWFpS7YN5JiAUXIv4Wer9xBjAJbGjGYmHQHRLA_UTuXgiAhlYC59QTdCmrmYq6f6CI-4IG10taVHAOmj9JJ9zr3rAB6IvXZ0ILhgzfaR4LbtfMn2-_zo7UyvLFeYabt-M-Tv03RCwbG4Whdbz5E92PNchVwCYSo7UHvW7xnGjkCtN5Tbcs9Zv1F86X5YG0WVqjpf7IhjFlTzaLdwWrjrcCNFGn_xLHzvIFIVfYAWSSSq1Arl5vTvER-0JQoYFRCvYVBexxK9SL-5eqsyQ_7AOwc3BS9tW8LJRoAI-CGkLP_MokTXHRBmHgijezvkzzM.DpZvBbRznHU52VLrzt9rmw",
			standCommonOutput)
		suite.Error(errParse, "token is expired")
	})
}

func (suite *JWERS256Suite) TestERS256JWTParseTokenMethodExpireNoExpired() {
	suite.NotPanics(func() {
		j, err := NewERS256JWT(suite.key)
		suite.Equal(nil, err)

		standCommonOutput := NewCommon(NewClaimsBuilder().Build())
		errParse := j.VerifyToken(
			"eyJhbGciOiJSU0EtT0FFUC0yNTYiLCJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2R0NNIiwidHlwIjoiSldUIn0.osB8ZnGtq7sB6Rc5bEPTDSCj3nO2tWb6MHWW3lta_gfImb9hlOZO6Mh6TegGgPrT0rAW7d82jEox607sA-4hSJolBid5YKeg6OH4U_qz5ieCsgFR6BKSeZO4ssinntOd0F-8dc8COOSw-US4n7RrvINWXyJPjVbKd298-NVlkCVl16iWX4WUwihD-SIlRyuTPmqufS3-EERZL_2WX17BAUrld91Q-LmbG9PcNDmNFwung7xMEb7E9EKQ5i9yucY9H2tTND9QW1oE5ZJJN4ZkJvUZy4AKQsMOqWwDxAwIvuTI9kNgKxAlENq_fTwL7-El801nrvrIFEE6ItZFWSzfNw.DegYBjJRI-_N_s8s.f90IxGVXKoZpEVIq-ubR3MDtX1wBXBLS4ZgGaWyUbn_VRRT0e_YVoZ1XHRgl3QBJwaEw9SidBWyI2KDlGrGUvKle4uWQOQNdmhQFOiubRm0FKfsvjpTZTnR_Q8EUne7L9OeqA92NTCxun72NkSvSDCsd2-h4elMH-1xj_XIAmfZeuI1pLVx6zBuDNYyPT7avSq27TQsJW89ls-8Sez52lf89vI2mWdo_avYaA3NHjPcADwiV1HNFcmUWyu6dxE-QUv6GPC60YaQyazbAaeOh2Ftc671ZI8VzCwGESXH5z8b3qwexuxTHIegMPOIhUAlc6SkMt6Lw5CqYt73v911ISltyHJP0fcQ8rqoZQsxOwAFw8R19NTi1XchT7vdAusU2J5gEGRjzsQXg9CLgiErBvPFNHsraem5UBZ-98vyKfG4u1mWJ0rtb6l9I5Z3bi2DV6UFbwkq4FdjL4_hLNsKxJWufmhpMJrgtNHRBdfPsVDgLAiKA3gHv2IHd1yq0iBavfFxoW2HHNZuqi3gCcTlA8GxxIfGwBPfs0MjJOMR6U-TdcqUYVGosX0K5o5_1zqdSKmcI1NRUD8MSYA8f2SL9Pg-fhh42KtwknoR_3Hx4iLsb0-Y5eWCLTWPfv3XIUbfxapru3SDJ996fxy757BQZ.BKwDPMP4O5sjEtt0xoKy3Q",
			standCommonOutput)
		suite.NoError(errParse)
	})
}

func (suite *JWERS256Suite) TestERS256JWTRefreshTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewERS256JWT(suite.key)
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
		newTk, errRefresh := j.RefreshToken(
			tk,
			standCommonOutput,
			100*time.Millisecond)
		suite.Equal(nil, errRefresh)
		suite.Equal(tk, newTk)
	})
}

func (suite *JWERS256Suite) TestERS256JWTRefreshTokenMethodExpire() {
	suite.NotPanics(func() {
		j, err := NewERS256JWT(suite.key)
		suite.Equal(nil, err)
		standClaimsOutput := NewClaimsBuilder().Build()
		tk, errRefresh := j.RefreshToken(
			"eyJhbGciOiJSU0EtT0FFUC0yNTYiLCJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2R0NNIiwidHlwIjoiSldUIn0.RPaJNb5RtS1CJQrSwrgw0185pJ2_RaZi0SdqlCSdiK_I7cprj4HOvx9TnIMASqxRRsdwYIZZN8reaYcYk53i27QsP7_2cfgOrATBeccDRoaNI6l-2twTrff1A-E-GiMJ3Po7ExdjIeIdB-dwmNHx2-r31AJJWURuNYEXRbNMmbX6eMMQQYNrz0KO_vmphUlS8b3vWhNrRORsKGkllmeEol_wQ8l7TVQ0DN-rLbwBM2J20enKwVV6fIZKeCQJLcKDBpj1Tz2J8GsbC5ks4YxneEE44Eq8D3tFJ0-7h_uQG4jY3yAR7WhFvrO0nCTF8pK-RrjVh4SQsJcRMMGJouqANQ.LN1Eh1ngtPsMTImQ.0qHmRSKIfQYmJiVzooIui0n5qjy0rCD24VBanXQAkxiy33LM-ur-3_fqOa_3hjTzn1L2JMQwtfS-IZNvQzjZyLlfgK8XlqV7iBDxzIHw7ewTAkPfC9OnQIzf8JolkWMUwknBezPl1_wkwR9H46kqKjpIIxprMOlDBjhTihKCxRyUQBQF6nt7yLZ7pSblMAF03X1c3P6Veoh2nC-8RNA-UXIQoq6cpmoDjeEsxqUhH0eqZIX9F5iHFT0LYE05STPIOZs3gSUPu9g1ZwTKGhaxeoKkNjVqVyheqdGDqJDitbLgIVygLzh-KfydciafqSoajb9xFsqrvw2M1XxMJwzpbuOOpmV_rOUvnToe_ZcQLvcJOPGa_fjqoM8OOqes52wfwA_WzeuvR_2k4danQEbqR1uI0Kb6aIPWmEf_2djI4lthdKn8qDYRgINfvjXQX7d9UY03zRT9fF5OYxvB8bw0WwUT04E9eQ-K4FPgix4OhZVzp0SIkq4euH_DdzT_i2rTU2tzcyqOYPOPkBcVWEAOTFeYWqQVBW1T6pqDT5RilpH5eqGTZ7CNJabX0LDyu8lXs_4ppudf4fsH5h52HVAzEz4JN0XIq4RxuqQJ6fLHr5thd2LDHf_fuRSDagJ4r7uGKsy2Tcqc1RmD1uZgFxM1xNBB3xgusIy4HCCnx8VfvhLa05SAlzw.1mrdAJhNbxGJP6ejA0zv5g",
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

func (suite *JWERS256Suite) TestERS256JWTRefreshTokenMethodparseRawError() {
	defer gostub.StubFunc(&parseSignedAndEncrypted, nil, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		j, err := NewERS256JWT(suite.key)
		suite.Equal(nil, err)
		standClaimsOutput := NewClaimsBuilder().Build()
		_, errRefresh := j.RefreshToken(
			"eyJhbGciOiJSU0EtT0FFUC0yNTYiLCJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2R0NNIiwidHlwIjoiSldUIn0.W81uCiTrZ1j-w311Q512RshQFpAIlXxlGnhlQVITwewfY8DoF0yJJgXQ8ifCvCu7n-IvD2cB7K2fmfiiCcJTsbKipHJg0DT9MxawV_JhaBRrzkKtzR842G4s8Z9EadefMcUFdaZ_klQ7I7P1ieMjDfd-HYKzleOip3BRRe0Y33AP803X5K584pTTUW9_vR2Vc305R1IF4qx9OpJiSZmNEZ-wP23WVhsegkBruDgPf0s4yJsVV1rdUZ9JFR9XNjq_6VvZ4klo1ODHqIZ1-4V3Gs1y8ShQKx2cb29gjoqtV01-7ik-CFWYmL83c6oTv_mjHPyCboVfrIAtIra1sJqOPA.tqaK3VNvv5_IZuwF.sdnCs6ovZICVQ50ypywZgEYeS7pgPHxdYoIBIg0Wm10LEpPCzcYT86ToufVzR4_IHGeHGZa7fiLGUJy4IzOvklhSjNGmqwWnZf5jfXQBwgizrTT0ZsdqnEHXEQWN4ATOoQ4kAnJRoh9gfJoOObPAilCWgvfZjE9Kq3cYNjbY7SJPhujmSLT5jwWiHx4-TbTGVU6CRmGUe5_EnfDvGtsxf0E1Pndevcu1aW5eOn6qYTFHX1H2NImuMBUc2sradY6KNPm_udd_qEs3jCCsGoDvtlLv3yEVhpkx1bHWiQdHx-2ZKnd4o8lI2jzjwojbae3VNqXaCaHjjcI2crpeYGKLlQ-xJpIVMRo3CCJZunoO_0dL6qLEIyeB_W7QfnWQBcddV1K8r5bbCU0AJbz6AqMNci51DO2_zCQEXVSMJCAAhdp8jo39leb8kfcgOkllZxi1P3xUexsy5r8ZwmuoB6s9Ou1-k1-UPUUen2XJhCtfB7-Hrz5BoxyrWsbtfLF22G3QAw-HwYltyrMbqpp8monT10usPZ2nJpp9F3H5AzdyGKphcJLB1uUll-c1Iu8RIGfcrch_ZuWAYDvn9QdpXFJNE_pgN6fUGYg7YXKBmYyNUlOZRsNrMMrg3x5L5q8WFGCecjZWjpXjjIeB-bnPISpH1zx22_stSsbjuIRjQA1F8-v75wvReO4.ns0TF_uNBL0HKBmvM2tEvQ",
			standClaimsOutput,
			100*time.Millisecond,
		)
		suite.Error(errRefresh, "got error")
	})
}

func TestJWERS256Suite(t *testing.T) {
	suite.Run(t, new(JWERS256Suite))
}
