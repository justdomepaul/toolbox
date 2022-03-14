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

type JWEHS256Suite struct {
	suite.Suite
	key    string
	option config.JWT
}

func (suite *JWEHS256Suite) SetupTest() {
	suite.key = `b583ed184e2018b3d89a4fa8832d0a1f`
	result := config.JWT{}
	suite.NoError(config.LoadFromEnv(&result))
	suite.option = result
}

func (suite *JWEHS256Suite) TestNewEHS256JWT() {
	suite.NotPanics(func() {
		_, err := NewEHS256JWT(suite.key)
		suite.Equal(nil, err)
	})
}

func (suite *JWEHS256Suite) TestNewEHS256JWTFromOptions() {
	suite.NotPanics(func() {
		_, err := NewEHS256JWTFromOptions(suite.option)
		suite.NoError(err)
	})
}

func (suite *JWEHS256Suite) TestNewEHS256JWTFromOptionsAllNoOption() {
	suite.NotPanics(func() {
		_, err := NewEHS256JWTFromOptions(config.JWT{
			HmacSecretKeyPath: "",
			HmacSecretKey:     "",
		})
		suite.Error(err, ErrNoKey)
	})
}

func (suite *JWEHS256Suite) TestNewEHS256JWTFromOptionsAllNoOptionNoFIle() {
	suite.NotPanics(func() {
		_, err := NewEHS256JWTFromOptions(config.JWT{
			HmacSecretKeyPath: "testFile.txt",
			HmacSecretKey:     "",
		})
		suite.Error(err)
	})
}

func (suite *JWEHS256Suite) TestNewEHS256JWTnewSignerError() {
	tSigner := testSigner{}
	defer gostub.StubFunc(&newSigner, tSigner, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewEHS256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWEHS256Suite) TestNewEHS256JWTnewEncrypterError() {
	tEncrypter := testEncrypter{}
	defer gostub.StubFunc(&newEncrypter, tEncrypter, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewEHS256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWEHS256Suite) TestEHS256JWTGenerateTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewEHS256JWT(suite.key)
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

func (suite *JWEHS256Suite) TestEHS256JWTValidateMethod() {
	suite.NotPanics(func() {
		j, err := NewEHS256JWT(suite.key)
		suite.Equal(nil, err)

		suite.NoError(j.Validate(
			"eyJhbGciOiJQQkVTMi1IUzUxMitBMjU2S1ciLCJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2Q0JDLUhTNTEyIiwicDJjIjoxMDAwMDAsInAycyI6IjFnWTIzZnhrVWloYnd3UUlSMTc1WkEiLCJ0eXAiOiJKV1QifQ.UOHbM5vuB5jZPs-PoseMtakFmwTj6YRQ6zWZ_9cYiYmJor9qLDtFy5rm767Jt_pQ8rr9UBRIETXgMEhV0Chy_PmYJ7NvRooW.pltiRrgShkSJNSP2PBZEQw.tn9t8Seg9qvMMZtgSPW6ENL--_tUgf22nYRa8RwpjUkxCNrRFPv1QNjZmykQReeadEgJVx_6olNR2aO0iGkvKp6O-bond7cw9YmGWxBnR0z9_mOWXjSYrZoz1AB3fco3z8JX5qsCOup2XsXbEuRX6nwLmt9NEEIrToj7Ae9lewBOZqUwAX85HbyXp0TcdkvgMTmHv_Ejud1mWDMa0Tt9_8Bztpe9g7xakm8fop6wGvhqL9TgfjN2vQFa0raMHzH2Q0qSvETvxofIc9IGqD6_TusTIThu6U4NYoS9s9DJTt6vfvZtFq7PzegSTpQFi-48wWuH3EwDJ90CFIb9YkSzaJp2DC1G15dwiICaMLlSIOt1PcVpleKirpjKQFA2VgAj.5KNs6Zq2AV4z1mBusvvwHdwOr7WGKqgAVgZx5gxnnpI"))
	})
}

func (suite *JWEHS256Suite) TestEHS256JWTValidateMethodparseSignedAndEncryptedError() {
	suite.NotPanics(func() {
		defer gostub.StubFunc(&parseSignedAndEncrypted, nil, errors.New("got error")).Reset()
		j, err := NewEHS256JWT(suite.key)
		suite.Equal(nil, err)

		suite.Error(j.Validate(
			"eyJhbGciOiJQQkVTMi1IUzUxMitBMjU2S1ciLCJjdHkiOiJKV1QiLCJlbmMiOiJBMjU2Q0JDLUhTNTEyIiwicDJjIjoxMDAwMDAsInAycyI6IjFnWTIzZnhrVWloYnd3UUlSMTc1WkEiLCJ0eXAiOiJKV1QifQ.UOHbM5vuB5jZPs-PoseMtakFmwTj6YRQ6zWZ_9cYiYmJor9qLDtFy5rm767Jt_pQ8rr9UBRIETXgMEhV0Chy_PmYJ7NvRooW.pltiRrgShkSJNSP2PBZEQw.tn9t8Seg9qvMMZtgSPW6ENL--_tUgf22nYRa8RwpjUkxCNrRFPv1QNjZmykQReeadEgJVx_6olNR2aO0iGkvKp6O-bond7cw9YmGWxBnR0z9_mOWXjSYrZoz1AB3fco3z8JX5qsCOup2XsXbEuRX6nwLmt9NEEIrToj7Ae9lewBOZqUwAX85HbyXp0TcdkvgMTmHv_Ejud1mWDMa0Tt9_8Bztpe9g7xakm8fop6wGvhqL9TgfjN2vQFa0raMHzH2Q0qSvETvxofIc9IGqD6_TusTIThu6U4NYoS9s9DJTt6vfvZtFq7PzegSTpQFi-48wWuH3EwDJ90CFIb9YkSzaJp2DC1G15dwiICaMLlSIOt1PcVpleKirpjKQFA2VgAj.5KNs6Zq2AV4z1mBusvvwHdwOr7WGKqgAVgZx5gxnnpI"))
	})
}

func (suite *JWEHS256Suite) TestEHS256JWTParseTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewEHS256JWT(suite.key)
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

func (suite *JWEHS256Suite) TestEHS256JWTParseTokenMethodExpire() {
	suite.NotPanics(func() {
		j, err := NewEHS256JWT(suite.key)
		suite.Equal(nil, err)

		standCommonOutput := NewCommon(NewClaimsBuilder().Build())
		errParse := j.VerifyToken(
			"eyJhbGciOiJQQkVTMi1IUzI1NitBMTI4S1ciLCJjdHkiOiJKV1QiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2IiwicDJjIjoxMDAwMDAsInAycyI6Ikw2Z3k3RXlJVkFZR1E0ZEJOQXN5VHciLCJ0eXAiOiJKV1QifQ.83YhoYO671b986kohh_xh_9cmcLJ9BS49JZwnbpgMlL4pWYcIbJ43w.ruMj8TmixURmR6htFYuL-w.Tmi48ifU_N7PWtNYDlzX6aJ8qMxHK_zyD3TQhZckfiVUOYVxGI_FApMsGjlA8QryIngP254F_6APpCftkAOSpZlt7VLMYptwYYjw-dHPE9n39RGSOYFAECKnrd0arox6FiroGJMJ9WyqB7vEibndQ68rS7S22G_c7pUVe3EJutaOnBYjeeR3oCTZDkfk-1KBv-GmavKTUYQhstBLLQdzmDWuLOBy-RQ7WUsO-wfw7hoxuCRru1R3FvRNF-fMCXCosMqtR6Axww0yu2dn3ya36g70x9hDz5SOrRERcDfrmvQ.RXg7ZgxIIMfNIKUQBXyPJw",
			standCommonOutput)
		suite.Error(errParse, "token is expired")
	})
}

func (suite *JWEHS256Suite) TestEHS256JWTParseTokenMethodExpireNoExpiredTime() {
	suite.NotPanics(func() {
		j, err := NewEHS256JWT(suite.key)
		suite.Equal(nil, err)

		standCommonOutput := NewCommon(NewClaimsBuilder().Build())
		errParse := j.VerifyToken(
			"eyJhbGciOiJQQkVTMi1IUzI1NitBMTI4S1ciLCJjdHkiOiJKV1QiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2IiwicDJjIjoxMDAwMDAsInAycyI6IktJU3Q1d2ZEZXBWQkw3c2VoQlcwS1EiLCJ0eXAiOiJKV1QifQ.VUnz_jxXUQAvIOrsiYxMsqxYT-Xl1-MSse0g6_3wtfRgN8gUBkRi5w.qR1SzUzPlhR9yQ_WbWeHDA.behikD4QiLk1H6_Ak8Ow3WHpiXpGHtG1Gpma_w6k9yOSr8UfeshEInHCIYZ9OMz88n-vx1AgRZ8mFZG0IdmALYzHA8uX4aD_y4nadNRv4NUQ73JnsbqUISEHHsqIW0dSZ8Qp8PBNuAN7eOi4divajfgyb2y2RP9Epvk9qs6TuJH-QJuncVhcY0hktA75JaALebEgcjJ8xG7FQuPXta838aZ_-S8BYPqKYW9WwA5yPnWXpb55lMGlcGoUEJ6-ahBmbl0R9uE3xz_MIr-SKO64DStBMxGS5Gsv_xGceAUNWaU.b_I-k-l8r6wYtuJsS9JN7g",
			standCommonOutput)
		suite.NoError(errParse)
	})
}

func (suite *JWEHS256Suite) TestEHS256JWTRefreshTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewEHS256JWT(suite.key)
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
		newTk, errRefresh := j.RefreshToken(
			tk,
			standCommonOutput,
			100*time.Millisecond)
		suite.Equal(nil, errRefresh)
		suite.Equal(tk, newTk)
	})
}

func (suite *JWEHS256Suite) TestEHS256JWTRefreshTokenMethodExpire() {
	suite.NotPanics(func() {
		j, err := NewEHS256JWT(suite.key)
		suite.Equal(nil, err)
		standClaimsOutput := NewClaimsBuilder().Build()
		tk, errRefresh := j.RefreshToken(
			"eyJhbGciOiJQQkVTMi1IUzI1NitBMTI4S1ciLCJjdHkiOiJKV1QiLCJlbmMiOiJBMTI4Q0JDLUhTMjU2IiwicDJjIjoxMDAwMDAsInAycyI6Ik4tSHAtdHMxQmdTYlU4M2NBazJWd1EiLCJ0eXAiOiJKV1QifQ.bEQR8CEbDIZI6Amz_beQ8FoWXp4aFfMIeOaPbTPNIacP6s4FZfMeiA.R8zE-i1kWtp9JPyfJeVHLQ.HvSSaLupbJDfU4mW-7dMEOz4ArqzoblTC8bxSE402KO_-ScQgue5-vXbSyUi_S4I3kags7S_0J-Icr1tBuL_Vplj2XYdn1HtM2Zb-HETabYAvHaDmJJUlgTR-ll3HCMSQwkaabtkYV6ZxZn4aib3Fi5K5gZ3dH1kWGMNwcEygrCK89Wq7lpK0ArE960hJn-usrjV0QoiJTVdVDcOJblqyzrbncoQq13qMJgjby2ob_ZeaPMW00tiyKzEo9dPbk77VL_8cKtsg1aW0yShDtd8mOArjDrNycc4PpU18t8SF0M.Ia7ljuFleGLV4ak3YSduSA",
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

func (suite *JWEHS256Suite) TestEHS256JWTRefreshTokenMethodparseRawError() {
	defer gostub.StubFunc(&parseSignedAndEncrypted, nil, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		j, err := NewEHS256JWT(suite.key)
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

func TestJWEHS256Suite(t *testing.T) {
	suite.Run(t, new(JWEHS256Suite))
}
