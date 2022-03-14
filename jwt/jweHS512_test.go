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

type JWEHS512Suite struct {
	suite.Suite
	key    string
	option config.JWT
}

func (suite *JWEHS512Suite) SetupTest() {
	suite.key = `b583ed184e2018b3d89a4fa8832d0a1f`
	result := config.JWT{}
	suite.NoError(config.LoadFromEnv(&result))
	suite.option = result
}

func (suite *JWEHS512Suite) TestNewEHS512JWT() {
	suite.NotPanics(func() {
		_, err := NewEHS512JWT(suite.key)
		suite.Equal(nil, err)
	})
}

func (suite *JWEHS512Suite) TestNewEHS512JWTFromOptions() {
	suite.NotPanics(func() {
		_, err := NewEHS512JWTFromOptions(suite.option)
		suite.NoError(err)
	})
}

func (suite *JWEHS512Suite) TestNewEHS512JWTFromOptionsAllNoOption() {
	suite.NotPanics(func() {
		_, err := NewEHS512JWTFromOptions(config.JWT{
			HmacSecretKeyPath: "",
			HmacSecretKey:     "",
		})
		suite.Error(err, ErrNoKey)
	})
}

func (suite *JWEHS512Suite) TestNewEHS512JWTFromOptionsAllNoOptionNoFIle() {
	suite.NotPanics(func() {
		_, err := NewEHS512JWTFromOptions(config.JWT{
			HmacSecretKeyPath: "testFile.txt",
			HmacSecretKey:     "",
		})
		suite.Error(err)
	})
}

func (suite *JWEHS512Suite) TestNewEHS512JWTnewSignerError() {
	tSigner := testSigner{}
	defer gostub.StubFunc(&newSigner, tSigner, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewEHS512JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWEHS512Suite) TestNewEHS512JWTnewEncrypterError() {
	tEncrypter := testEncrypter{}
	defer gostub.StubFunc(&newEncrypter, tEncrypter, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewEHS512JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWEHS512Suite) TestEHS512JWTGenerateTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewEHS512JWT(suite.key)
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

func (suite *JWEHS512Suite) TestEHS512JWTValidateMethod() {
	suite.NotPanics(func() {
		j, err := NewEHS512JWT(suite.key)
		suite.Equal(nil, err)

		suite.NoError(j.Validate(
			"eyJhbGciOiJQQkVTMi1IUzUxMitBMjU2S1ciLCJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2Q0JDLUhTNTEyIiwicDJjIjoxMDAwMDAsInAycyI6IjFnWTIzZnhrVWloYnd3UUlSMTc1WkEiLCJ0eXAiOiJKV1QifQ.UOHbM5vuB5jZPs-PoseMtakFmwTj6YRQ6zWZ_9cYiYmJor9qLDtFy5rm767Jt_pQ8rr9UBRIETXgMEhV0Chy_PmYJ7NvRooW.pltiRrgShkSJNSP2PBZEQw.tn9t8Seg9qvMMZtgSPW6ENL--_tUgf22nYRa8RwpjUkxCNrRFPv1QNjZmykQReeadEgJVx_6olNR2aO0iGkvKp6O-bond7cw9YmGWxBnR0z9_mOWXjSYrZoz1AB3fco3z8JX5qsCOup2XsXbEuRX6nwLmt9NEEIrToj7Ae9lewBOZqUwAX85HbyXp0TcdkvgMTmHv_Ejud1mWDMa0Tt9_8Bztpe9g7xakm8fop6wGvhqL9TgfjN2vQFa0raMHzH2Q0qSvETvxofIc9IGqD6_TusTIThu6U4NYoS9s9DJTt6vfvZtFq7PzegSTpQFi-48wWuH3EwDJ90CFIb9YkSzaJp2DC1G15dwiICaMLlSIOt1PcVpleKirpjKQFA2VgAj.5KNs6Zq2AV4z1mBusvvwHdwOr7WGKqgAVgZx5gxnnpI"))
	})
}

func (suite *JWEHS512Suite) TestEHS512JWTValidateMethodparseSignedAndEncryptedError() {
	suite.NotPanics(func() {
		defer gostub.StubFunc(&parseSignedAndEncrypted, nil, errors.New("got error")).Reset()
		j, err := NewEHS512JWT(suite.key)
		suite.Equal(nil, err)

		suite.Error(j.Validate(
			"eyJhbGciOiJQQkVTMi1IUzUxMitBMjU2S1ciLCJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2Q0JDLUhTNTEyIiwicDJjIjoxMDAwMDAsInAycyI6IjFnWTIzZnhrVWloYnd3UUlSMTc1WkEiLCJ0eXAiOiJKV1QifQ.UOHbM5vuB5jZPs-PoseMtakFmwTj6YRQ6zWZ_9cYiYmJor9qLDtFy5rm767Jt_pQ8rr9UBRIETXgMEhV0Chy_PmYJ7NvRooW.pltiRrgShkSJNSP2PBZEQw.tn9t8Seg9qvMMZtgSPW6ENL--_tUgf22nYRa8RwpjUkxCNrRFPv1QNjZmykQReeadEgJVx_6olNR2aO0iGkvKp6O-bond7cw9YmGWxBnR0z9_mOWXjSYrZoz1AB3fco3z8JX5qsCOup2XsXbEuRX6nwLmt9NEEIrToj7Ae9lewBOZqUwAX85HbyXp0TcdkvgMTmHv_Ejud1mWDMa0Tt9_8Bztpe9g7xakm8fop6wGvhqL9TgfjN2vQFa0raMHzH2Q0qSvETvxofIc9IGqD6_TusTIThu6U4NYoS9s9DJTt6vfvZtFq7PzegSTpQFi-48wWuH3EwDJ90CFIb9YkSzaJp2DC1G15dwiICaMLlSIOt1PcVpleKirpjKQFA2VgAj.5KNs6Zq2AV4z1mBusvvwHdwOr7WGKqgAVgZx5gxnnpI"))
	})
}

func (suite *JWEHS512Suite) TestEHS512JWTParseTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewEHS512JWT(suite.key)
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

func (suite *JWEHS512Suite) TestEHS512JWTParseTokenMethodExpire() {
	suite.NotPanics(func() {
		j, err := NewEHS512JWT(suite.key)
		suite.Equal(nil, err)

		standCommonOutput := NewCommon(NewClaimsBuilder().Build())
		errParse := j.VerifyToken(
			"eyJhbGciOiJQQkVTMi1IUzUxMitBMjU2S1ciLCJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2Q0JDLUhTNTEyIiwicDJjIjoxMDAwMDAsInAycyI6Im9IQjk0dXlPMHVRQmdnVjdUaDBiRWciLCJ0eXAiOiJKV1QifQ.tRJluGb5G-rgvyfRHv9jGo2dm2e3OoLgkg2Lk3wJHA5ztqO9uy5t1D_WeWVAkso26m7CcI3-1ZjvbgAwn9ouZKO2vFCedjHq.ichM0zvvUJ2i0CFOoNMubg.b1PcccdftHkuI3U3cCzrF9g9S0FRactHiKbn9GrlWcBv02cI8_EEADj3OMGOwZt4AOD2ic0R7HGabw9aRLrNwV62c0pfsFoOJas2SVPJubAKklnTR9SISzUogmO16quXfspuJ-bTenfBlWTozmvkemFgQaPpLxilBlyOChjW1zTIGDz5w-1R66WEr_NcUoBpPEKADHEkwZJA-LiSIalnFZLqw83Wi3H1wJz2glEkr4GFtS1sQ5C6lA0zVzxkG9W99jJEsuSEAAgdO0KihNwEXvT02yb9WslL85lwoyYHOVz0KqexQIrwmaw7zPPa7S1g9jVGMZFPTx6KyeBYO6hAD_R7uwpsu-Uf0SVdacHiXEU.HpFVwVGZlutvl8k1jL5BO8bkEofLioasGWoQ0O62qpg",
			standCommonOutput)
		suite.Error(errParse, "token is expired")
	})
}

func (suite *JWEHS512Suite) TestEHS512JWTParseTokenMethodExpireNoExpired() {
	suite.NotPanics(func() {
		j, err := NewEHS512JWT(suite.key)
		suite.Equal(nil, err)

		standCommonOutput := NewCommon(NewClaimsBuilder().Build())
		errParse := j.VerifyToken(
			"eyJhbGciOiJQQkVTMi1IUzUxMitBMjU2S1ciLCJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2Q0JDLUhTNTEyIiwicDJjIjoxMDAwMDAsInAycyI6ImhSbEpOX3AzcVJTSnpBYWxKN0MtUWciLCJ0eXAiOiJKV1QifQ.vlnMrChEbHoaXtWN0DgN55vfwKMnWUVHJmaaNgMqJhZwb9vVKJ5GpxGqbibiRwUeMwEaqRgcuR0wL939eiMn-ziJURx2s2n2.L3EwJPC9i6koAS1j_UyTNA.vH9AXKWCVk3Z3iumwutfRw1IYkXQRU5boFAz88bNCQDUz42w9cSmntTPryq15U_K1e0TqP-iHDqptIOOXvBHibQFprq8wYH_IYJdd-ZAbNhGUMpNpM1szWeCbljhDxx_rqq35DCqg7oWUOAQeUIjzAHrjtUAdFYQH7YhLNfQCk1NPSq3x7kynOwFL4zS5oGUEu-2sxcyazuWxNDV7AN3qIESPfl4FnUUkq_OWhbCUYU6zOBwKZNrT6IEuqaGzVu97XlnhS5coMqlPE3DvstsMAYJcpFQ6ywgkmrw9U2y5GBte45IctOQ4DZGaZ9mB-L78vA-alfPJHuC4c_nhPfMDA.6h9Xyd6KaI712Iavpe2WbTRA9S67OMZfPQbmqw1RA5U",
			standCommonOutput)
		suite.NoError(errParse)
	})
}

func (suite *JWEHS512Suite) TestEHS512JWTRefreshTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewEHS512JWT(suite.key)
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

func (suite *JWEHS512Suite) TestEHS512JWTRefreshTokenMethodExpire() {
	suite.NotPanics(func() {
		j, err := NewEHS512JWT(suite.key)
		suite.Equal(nil, err)
		standClaimsOutput := NewClaimsBuilder().Build()
		tk, errRefresh := j.RefreshToken(
			"eyJhbGciOiJQQkVTMi1IUzUxMitBMjU2S1ciLCJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2Q0JDLUhTNTEyIiwicDJjIjoxMDAwMDAsInAycyI6Im9IQjk0dXlPMHVRQmdnVjdUaDBiRWciLCJ0eXAiOiJKV1QifQ.tRJluGb5G-rgvyfRHv9jGo2dm2e3OoLgkg2Lk3wJHA5ztqO9uy5t1D_WeWVAkso26m7CcI3-1ZjvbgAwn9ouZKO2vFCedjHq.ichM0zvvUJ2i0CFOoNMubg.b1PcccdftHkuI3U3cCzrF9g9S0FRactHiKbn9GrlWcBv02cI8_EEADj3OMGOwZt4AOD2ic0R7HGabw9aRLrNwV62c0pfsFoOJas2SVPJubAKklnTR9SISzUogmO16quXfspuJ-bTenfBlWTozmvkemFgQaPpLxilBlyOChjW1zTIGDz5w-1R66WEr_NcUoBpPEKADHEkwZJA-LiSIalnFZLqw83Wi3H1wJz2glEkr4GFtS1sQ5C6lA0zVzxkG9W99jJEsuSEAAgdO0KihNwEXvT02yb9WslL85lwoyYHOVz0KqexQIrwmaw7zPPa7S1g9jVGMZFPTx6KyeBYO6hAD_R7uwpsu-Uf0SVdacHiXEU.HpFVwVGZlutvl8k1jL5BO8bkEofLioasGWoQ0O62qpg",
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

func (suite *JWEHS512Suite) TestEHS512JWTRefreshTokenMethodparseRawError() {
	defer gostub.StubFunc(&parseSignedAndEncrypted, nil, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		j, err := NewEHS512JWT(suite.key)
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

func TestJWEHS512Suite(t *testing.T) {
	suite.Run(t, new(JWEHS512Suite))
}
