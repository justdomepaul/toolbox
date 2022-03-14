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

type JWEHS384Suite struct {
	suite.Suite
	key    string
	option config.JWT
}

func (suite *JWEHS384Suite) SetupTest() {
	suite.key = `b583ed184e2018b3d89a4fa8832d0a1f`
	result := config.JWT{}
	suite.NoError(config.LoadFromEnv(&result))
	suite.option = result
}

func (suite *JWEHS384Suite) TestNewEHS384JWT() {
	suite.NotPanics(func() {
		_, err := NewEHS384JWT(suite.key)
		suite.Equal(nil, err)
	})
}

func (suite *JWEHS384Suite) TestNewEHS384JWTFromOptions() {
	suite.NotPanics(func() {
		_, err := NewEHS384JWTFromOptions(suite.option)
		suite.NoError(err)
	})
}

func (suite *JWEHS384Suite) TestNewEHS384JWTFromOptionsAllNoOption() {
	suite.NotPanics(func() {
		_, err := NewEHS384JWTFromOptions(config.JWT{
			HmacSecretKeyPath: "",
			HmacSecretKey:     "",
		})
		suite.Error(err, ErrNoKey)
	})
}

func (suite *JWEHS384Suite) TestNewEHS384JWTFromOptionsAllNoOptionNoFIle() {
	suite.NotPanics(func() {
		_, err := NewEHS384JWTFromOptions(config.JWT{
			HmacSecretKeyPath: "testFile.txt",
			HmacSecretKey:     "",
		})
		suite.Error(err)
	})
}

func (suite *JWEHS384Suite) TestNewEHS384JWTnewSignerError() {
	tSigner := testSigner{}
	defer gostub.StubFunc(&newSigner, tSigner, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewEHS384JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWEHS384Suite) TestNewEHS384JWTnewEncrypterError() {
	tEncrypter := testEncrypter{}
	defer gostub.StubFunc(&newEncrypter, tEncrypter, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewEHS384JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWEHS384Suite) TestEHS384JWTGenerateTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewEHS384JWT(suite.key)
		suite.Equal(nil, err)
		standClaims := NewClaimsBuilder().
			WithSubject("testTopic").
			WithIssuer("tester").
			WithID("test001").
			WithAudience([]string{"testerClient"}).
			ExpiresAfter(5 * time.Second)
		//ExpiresAfter(87600 * time.Hour)
		common := NewCommon(standClaims.Build(), WithSecret("testData"))
		tk, errTk := j.GenerateToken(common)
		suite.T().Log(tk)
		suite.Equal(nil, errTk)
		suite.NotEmpty(tk)
	})
}

func (suite *JWEHS384Suite) TestEHS384JWTValidateMethod() {
	suite.NotPanics(func() {
		j, err := NewEHS384JWT(suite.key)
		suite.Equal(nil, err)

		suite.NoError(j.Validate(
			"eyJhbGciOiJQQkVTMi1IUzUxMitBMjU2S1ciLCJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2Q0JDLUhTNTEyIiwicDJjIjoxMDAwMDAsInAycyI6IjFnWTIzZnhrVWloYnd3UUlSMTc1WkEiLCJ0eXAiOiJKV1QifQ.UOHbM5vuB5jZPs-PoseMtakFmwTj6YRQ6zWZ_9cYiYmJor9qLDtFy5rm767Jt_pQ8rr9UBRIETXgMEhV0Chy_PmYJ7NvRooW.pltiRrgShkSJNSP2PBZEQw.tn9t8Seg9qvMMZtgSPW6ENL--_tUgf22nYRa8RwpjUkxCNrRFPv1QNjZmykQReeadEgJVx_6olNR2aO0iGkvKp6O-bond7cw9YmGWxBnR0z9_mOWXjSYrZoz1AB3fco3z8JX5qsCOup2XsXbEuRX6nwLmt9NEEIrToj7Ae9lewBOZqUwAX85HbyXp0TcdkvgMTmHv_Ejud1mWDMa0Tt9_8Bztpe9g7xakm8fop6wGvhqL9TgfjN2vQFa0raMHzH2Q0qSvETvxofIc9IGqD6_TusTIThu6U4NYoS9s9DJTt6vfvZtFq7PzegSTpQFi-48wWuH3EwDJ90CFIb9YkSzaJp2DC1G15dwiICaMLlSIOt1PcVpleKirpjKQFA2VgAj.5KNs6Zq2AV4z1mBusvvwHdwOr7WGKqgAVgZx5gxnnpI"))
	})
}

func (suite *JWEHS384Suite) TestEHS384JWTValidateMethodparseSignedAndEncryptedError() {
	suite.NotPanics(func() {
		defer gostub.StubFunc(&parseSignedAndEncrypted, nil, errors.New("got error")).Reset()
		j, err := NewEHS384JWT(suite.key)
		suite.Equal(nil, err)

		suite.Error(j.Validate(
			"eyJhbGciOiJQQkVTMi1IUzUxMitBMjU2S1ciLCJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2Q0JDLUhTNTEyIiwicDJjIjoxMDAwMDAsInAycyI6IjFnWTIzZnhrVWloYnd3UUlSMTc1WkEiLCJ0eXAiOiJKV1QifQ.UOHbM5vuB5jZPs-PoseMtakFmwTj6YRQ6zWZ_9cYiYmJor9qLDtFy5rm767Jt_pQ8rr9UBRIETXgMEhV0Chy_PmYJ7NvRooW.pltiRrgShkSJNSP2PBZEQw.tn9t8Seg9qvMMZtgSPW6ENL--_tUgf22nYRa8RwpjUkxCNrRFPv1QNjZmykQReeadEgJVx_6olNR2aO0iGkvKp6O-bond7cw9YmGWxBnR0z9_mOWXjSYrZoz1AB3fco3z8JX5qsCOup2XsXbEuRX6nwLmt9NEEIrToj7Ae9lewBOZqUwAX85HbyXp0TcdkvgMTmHv_Ejud1mWDMa0Tt9_8Bztpe9g7xakm8fop6wGvhqL9TgfjN2vQFa0raMHzH2Q0qSvETvxofIc9IGqD6_TusTIThu6U4NYoS9s9DJTt6vfvZtFq7PzegSTpQFi-48wWuH3EwDJ90CFIb9YkSzaJp2DC1G15dwiICaMLlSIOt1PcVpleKirpjKQFA2VgAj.5KNs6Zq2AV4z1mBusvvwHdwOr7WGKqgAVgZx5gxnnpI"))
	})
}

func (suite *JWEHS384Suite) TestEHS384JWTParseTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewEHS384JWT(suite.key)
		suite.Equal(nil, err)
		standClaims := NewClaimsBuilder().
			WithSubject("testTopic").
			WithIssuer("tester").
			WithID("test001").
			WithAudience([]string{"testerClient"}).
			ExpiresAfter(3 * time.Second)
		common := NewCommon(standClaims.Build(), WithSecret("testData"))
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

func (suite *JWEHS384Suite) TestEHS384JWTParseTokenMethodExpire() {
	suite.NotPanics(func() {
		j, err := NewEHS384JWT(suite.key)
		suite.Equal(nil, err)

		standCommonOutput := NewCommon(NewClaimsBuilder().Build())
		errParse := j.VerifyToken(
			"eyJhbGciOiJQQkVTMi1IUzM4NCtBMTkyS1ciLCJjdHkiOiJKV1QiLCJlbmMiOiJBMTkyQ0JDLUhTMzg0IiwicDJjIjoxMDAwMDAsInAycyI6IjUzRlZHc3J6MjBXVWNMMDdMOWtoaFEiLCJ0eXAiOiJKV1QifQ.5drt_UxBTw_Qenb9j_0ECQDhe1g7eNzDXHk2k6s_HkNxCKUM_AwGJqCtMCsCYS-pprhA1ZxFlW8.3Mb6-tCO2h8AiAi7jvuLpg.QLvMLp5WoYDUPCf3O-F8QzBgrrRHe0GxbW3N5IOLdAmmsNh-K5-O4x3VBwUHFpU9B4A8MTYkjc6AFV2z0_p4bsCbCSsLLgyErX_ucNMSLg9Zrx6LnnITfghW8wqkZuRZ4gcpEXhDiQahOcXuw14N3ounOo1NqrFVaG8ei9J5xeecIUTfR4C5iNVSAH8gkw_bcqe1bLP9woo04VsTnyeQ1Uc0hPyG9HBixjrWqWcqHPoDKWaz8bFXRGn_NOE1xaAI0PXI9mLzWiebT8xM71KGadPKYUYky6vrXXM48i2fZv7duh-xQsHp9Xd-rHf3kf4C2UZhjblCs2Q2LcBQ562dmQ.KpV7wvMNmDy6uc7lUonPbDivRrXbJe1l",
			standCommonOutput)
		suite.Error(errParse, "token is expired")
	})
}

func (suite *JWEHS384Suite) TestEHS384JWTParseTokenMethodExpireNoExpired() {
	suite.NotPanics(func() {
		j, err := NewEHS384JWT(suite.key)
		suite.Equal(nil, err)

		standCommonOutput := NewCommon(NewClaimsBuilder().Build())
		errParse := j.VerifyToken(
			"eyJhbGciOiJQQkVTMi1IUzM4NCtBMTkyS1ciLCJjdHkiOiJKV1QiLCJlbmMiOiJBMTkyQ0JDLUhTMzg0IiwicDJjIjoxMDAwMDAsInAycyI6ImdsSVYwdk5DLTJJUTMwRnZ5dHp3cnciLCJ0eXAiOiJKV1QifQ.9oOhhvgZNPVXY7rlJu9NS0-5vsneBV5uswIl2kL6S2pYV2zOxiFtWZZfungTW0Md1ReXk52lcMY.QM7XfuoXTAKZedeaQqhsPA.yFGe7MJUVs5ny3a0ToyLoy_goFl-tp7CucKp9MNunUFF20ysGvcLFPBn4re-Z3B2xhlNY2N5JgrXPOY8sMNueVX8yingnJKY1p43sDkGXFnNSetipg5QCMGfNhPNiPeO6Is9cB155gDhtHI2KnjL3SFczk8TTMSVFHfLJ4fPJcwdA_mRZiptXXQRDMxP6cWYja2olB2Cgk2Aeyo79HiHZJjMkdxxyFxymNx4nNLwwV4KYIqWf-0ueUGxZJW7v2sbkkhhp2WMSgfq3txKKEGrWP1JUR5ZqmQjU9Q9cMtTh6Br_ggLnN6ZLuz7yx5PYIwP.cP1M64hJ9BFkotYz_q630jWDu9KD8RpG",
			standCommonOutput)
		suite.NoError(errParse)
	})
}

func (suite *JWEHS384Suite) TestEHS384JWTRefreshTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewEHS384JWT(suite.key)
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

func (suite *JWEHS384Suite) TestEHS384JWTRefreshTokenMethodExpire() {
	suite.NotPanics(func() {
		j, err := NewEHS384JWT(suite.key)
		suite.Equal(nil, err)
		standClaimsOutput := NewClaimsBuilder().Build()
		tk, errRefresh := j.RefreshToken(
			"eyJhbGciOiJQQkVTMi1IUzM4NCtBMTkyS1ciLCJjdHkiOiJKV1QiLCJlbmMiOiJBMTkyQ0JDLUhTMzg0IiwicDJjIjoxMDAwMDAsInAycyI6ImYtaGdNQWQ3MEN4SUthRlBTSmFuYmciLCJ0eXAiOiJKV1QifQ.OoqdiGCfzuyaAdNovJudb_IGHa09VSHM5Rf2VTliYit6kDjhrDiz77Rv9gQoB8oMcZWphG9aqXQ.eIpL-HHdjg3RKvWvigH8Gg.tjPRnIlEeJ8DbLZccErmZrFLLhq6GE6ZWQMMRGwUUN_eet27z_RzJVgnaCO0_0ay7A1LBbYIG1ubmrLcKs_8YpQeDOAk_lpnT6qFoHgkXKlgq5uBS7kHjer71CUpRHm2SvshaHPCmUohpc7IBVmg-KR4X5aIvxjN1tQ8IlvL2KqXUlYLr72F1Ow_28ahnov7sLOqIjukDLEfcxj2xAanAP7DYZS8-6FLOmUjdUpgfKTPwj6nE4kyLneUCqvoA8y9K4NETO4udrUyfzI9nZfcmMAlndiMWIvpKh2uTtuzcJLmut4cgSVBzadYgDWplLTrUoFrtHcc3HT1gxQPhh11Hg.hFQGI64wRj_HPqNt1RGCW9MIz7hknxe5",
			standClaimsOutput,
			100*time.Millisecond,
		)
		suite.NoError(errRefresh)
		suite.NotEmpty(tk)
		standClaimsOutputOther := NewClaimsBuilder().Build()
		standCommonOutput := NewCommon(standClaimsOutputOther)
		time.Sleep(105 * time.Millisecond)
		errParse := j.VerifyToken(tk, standCommonOutput)
		suite.Error(errParse, ErrTokenExpired)
	})
}

func (suite *JWEHS384Suite) TestEHS384JWTRefreshTokenMethodparseRawError() {
	defer gostub.StubFunc(&parseSignedAndEncrypted, nil, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		j, err := NewEHS384JWT(suite.key)
		suite.Equal(nil, err)
		standClaimsOutput := NewClaimsBuilder().Build()
		_, errRefresh := j.RefreshToken(
			"eyJhbGciOiJQQkVTMi1IUzUxMitBMjU2S1ciLCJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2Q0JDLUhTNTEyIiwicDJjIjoxMDAwMDAsInAycyI6IjFnWTIzZnhrVWloYnd3UUlSMTc1WkEiLCJ0eXAiOiJKV1QifQ.UOHbM5vuB5jZPs-PoseMtakFmwTj6YRQ6zWZ_9cYiYmJor9qLDtFy5rm767Jt_pQ8rr9UBRIETXgMEhV0Chy_PmYJ7NvRooW.pltiRrgShkSJNSP2PBZEQw.tn9t8Seg9qvMMZtgSPW6ENL--_tUgf22nYRa8RwpjUkxCNrRFPv1QNjZmykQReeadEgJVx_6olNR2aO0iGkvKp6O-bond7cw9YmGWxBnR0z9_mOWXjSYrZoz1AB3fco3z8JX5qsCOup2XsXbEuRX6nwLmt9NEEIrToj7Ae9lewBOZqUwAX85HbyXp0TcdkvgMTmHv_Ejud1mWDMa0Tt9_8Bztpe9g7xakm8fop6wGvhqL9TgfjN2vQFa0raMHzH2Q0qSvETvxofIc9IGqD6_TusTIThu6U4NYoS9s9DJTt6vfvZtFq7PzegSTpQFi-48wWuH3EwDJ90CFIb9YkSzaJp2DC1G15dwiICaMLlSIOt1PcVpleKirpjKQFA2VgAj.5KNs6Zq2AV4z1mBusvvwHdwOr7WGKqgAVgZx5gxnnpI",
			standClaimsOutput,
			100*time.Millisecond,
		)
		suite.Error(errRefresh, "got error")
	})
}

func TestJWEHS384Suite(t *testing.T) {
	suite.Run(t, new(JWEHS384Suite))
}
