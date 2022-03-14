package jwt

import (
	"github.com/cockroachdb/errors"
	"github.com/justdomepaul/toolbox/config"
	"github.com/prashantv/gostub"
	"github.com/square/go-jose/v3"
	"github.com/square/go-jose/v3/jwt"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type testSigner struct {
	jose.Signer
}

type testEncrypter struct {
	jose.Encrypter
}

type JWEES256Suite struct {
	suite.Suite
	key    string
	option config.JWT
}

func (suite *JWEES256Suite) SetupTest() {
	suite.key = `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIANpsaiXN0UdToO728+5kwGNPA+RsbEPwb3MGiAlBhSXoAoGCCqGSM49
AwEHoUQDQgAER5WNvPs/SMICGESgDbN7IYl0CvPSkUhAaUtF/LAQEINqte/HLMks
hRsKJ2MTCe1upn5vhgBuGl5CL4ea4DqNhA==
-----END EC PRIVATE KEY-----
`
	result := config.JWT{}
	suite.NoError(config.LoadFromEnv(&result))
	suite.option = result
}

func (suite *JWEES256Suite) TestNewEES256JWT() {
	suite.NotPanics(func() {
		_, err := NewEES256JWT(suite.key)
		suite.NoError(err)
	})
}

func (suite *JWEES256Suite) TestNewEES256JWTFromOptions() {
	suite.NotPanics(func() {
		_, err := NewEES256JWTFromOptions(suite.option)
		suite.NoError(err)
	})
}

func (suite *JWEES256Suite) TestNewEES256JWTFromOptionsAllNoOption() {
	suite.NotPanics(func() {
		_, err := NewEES256JWTFromOptions(config.JWT{
			EcdsaPrivateKeyPath: "",
			EcdsaPrivateKey:     "",
		})
		suite.Error(err, ErrNoKey)
	})
}

func (suite *JWEES256Suite) TestNewEES256JWTFromOptionsAllNoOptionNoFIle() {
	suite.NotPanics(func() {
		_, err := NewEES256JWTFromOptions(config.JWT{
			EcdsaPrivateKeyPath: "testFile.txt",
			EcdsaPrivateKey:     "",
		})
		suite.Error(err)
	})
}

func (suite *JWEES256Suite) TestNewEES256JWTparseECPrivateKeyFromPEMError() {
	defer gostub.StubFunc(&parseECPrivateKeyFromPEM, nil, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewEES256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWEES256Suite) TestNewEES256JWTnewSignerError() {
	tSigner := testSigner{}
	defer gostub.StubFunc(&newSigner, tSigner, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewEES256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWEES256Suite) TestNewEES256JWTnewEncrypterError() {
	tEncrypter := testEncrypter{}
	defer gostub.StubFunc(&newEncrypter, tEncrypter, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewEES256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWEES256Suite) TestEES256JWTGenerateTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewEES256JWT(suite.key)
		suite.Equal(nil, err)
		standClaims := NewClaimsBuilder().
			WithSubject("testTopic").
			WithIssuer("tester").
			WithID("test001").
			WithAudience([]string{"testerClient"}).Build()
		//ExpiresAfter(87600 * time.Hour)
		common := NewCommon(standClaims, WithSecret("testData"))
		tk, errTk := j.GenerateToken(common)
		suite.T().Log(tk)
		suite.Equal(nil, errTk)
		suite.NotEmpty(tk)
	})
}

func (suite *JWEES256Suite) TestEES256JWTValidateMethod() {
	suite.NotPanics(func() {
		j, err := NewEES256JWT(suite.key)
		suite.Equal(nil, err)

		suite.NoError(j.Validate(
			"eyJhbGciOiJFQ0RILUVTK0EyNTZLVyIsImN0eSI6IkpXVCIsImVuYyI6IkEyNTZHQ00iLCJlcGsiOnsia3R5IjoiRUMiLCJjcnYiOiJQLTI1NiIsIngiOiJRZmJVRnNKM3NGV3ZKMkkxRGtQUElmYjQ3dGRTSmxsUzB3VjJ1VHhuQ2ZJIiwieSI6IjlXelNfTE5tZGkybzZUb3BVMXNwTnpfVlUwVk1QS3JZWkV4c0s5ODhRdHcifSwidHlwIjoiSldUIn0.AZwfbZgE2ib6ANat7lOVBBNfcE31Uccy_YlCedxkchL60xDxYAJqkA.lnzjMHoNKuNU57rZ.b6Vsp3g-ntbKH1J2XWkUrp9k7NCf-oNf6cTQ6lI800FAaplNIuRaZwQ_D7Rh3dUkSyNToYICgDKclNbVAivqTrqhbh6VJcUI7hbHMWIVCMu3Jh3TA6ZigwlrB9NZOZknvljpjLVMb7DZLZwC3rZnyinNHQTPISRbKO4Xqoe1g90_UwCEokG8UNm9GFsrudrfShunirj9xv42Z0aIKyDefrEsT3dQ3pGeX9uJl09S6RHYgEhLAYZ463_n7wixFf6OiPEuTIeQEp7nDi8HXGJkW3YIn0RbLK5V_RIvxCHgYX8SHlCYEzrAi97YjjplUdM._njneQ_kBKGFueqEjPdL5g"))
	})
}

func (suite *JWEES256Suite) TestEES256JWTValidateMethodparseSignedAndEncryptedError() {
	suite.NotPanics(func() {
		defer gostub.StubFunc(&parseSignedAndEncrypted, nil, errors.New("got error")).Reset()
		j, err := NewEES256JWT(suite.key)
		suite.Equal(nil, err)

		suite.Error(j.Validate(
			"eyJhbGciOiJFQ0RILUVTK0EyNTZLVyIsImN0eSI6IkpXVCIsImVuYyI6IkEyNTZHQ00iLCJlcGsiOnsia3R5IjoiRUMiLCJjcnYiOiJQLTI1NiIsIngiOiJRZmJVRnNKM3NGV3ZKMkkxRGtQUElmYjQ3dGRTSmxsUzB3VjJ1VHhuQ2ZJIiwieSI6IjlXelNfTE5tZGkybzZUb3BVMXNwTnpfVlUwVk1QS3JZWkV4c0s5ODhRdHcifSwidHlwIjoiSldUIn0.AZwfbZgE2ib6ANat7lOVBBNfcE31Uccy_YlCedxkchL60xDxYAJqkA.lnzjMHoNKuNU57rZ.b6Vsp3g-ntbKH1J2XWkUrp9k7NCf-oNf6cTQ6lI800FAaplNIuRaZwQ_D7Rh3dUkSyNToYICgDKclNbVAivqTrqhbh6VJcUI7hbHMWIVCMu3Jh3TA6ZigwlrB9NZOZknvljpjLVMb7DZLZwC3rZnyinNHQTPISRbKO4Xqoe1g90_UwCEokG8UNm9GFsrudrfShunirj9xv42Z0aIKyDefrEsT3dQ3pGeX9uJl09S6RHYgEhLAYZ463_n7wixFf6OiPEuTIeQEp7nDi8HXGJkW3YIn0RbLK5V_RIvxCHgYX8SHlCYEzrAi97YjjplUdM._njneQ_kBKGFueqEjPdL5g"))
	})
}

func (suite *JWEES256Suite) TestEES256JWTParseTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewEES256JWT(suite.key)
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

func (suite *JWEES256Suite) TestEES256JWTParseTokenMethodExpire() {
	suite.NotPanics(func() {
		j, err := NewEES256JWT(suite.key)
		suite.Equal(nil, err)

		standCommonOutput := NewCommon(NewClaimsBuilder().Build())
		errParse := j.VerifyToken(
			`eyJhbGciOiJFQ0RILUVTK0EyNTZLVyIsImN0eSI6IkpXVCIsImVuYyI6IkEyNTZHQ00iLCJlcGsiOnsia3R5IjoiRUMiLCJjcnYiOiJQLTI1NiIsIngiOiJwSGhWcVNlWG5NNExQWlB6ZkdZdEJXNW14WW10R1FITXRqbVJFTzhUaW9VIiwieSI6InlKdFN2cDFqTjJSUHJvVGJDY3VGbTZHeXdZanVrOV9MMzdZdkE2dWF6NXMifSwidHlwIjoiSldUIn0.pnNybB-cJ8oXEiqXCkA6-X_Qwu1VDhAxQz3Evy6c2wzXYxscaovJew.sP5P22e1dthpox5S.Tvc2E_sfFqyY8vDsg0LrU_Ob-S4em38MNi2cZurf8CVbEWd9j82NzDn6WIkwe2ZW-kKhWUrA7PHbAY7SjqNLgh1jUlTIahArF7S1wpN_uPm91SOzhAd54UMlh4aGfzvU8tZLsDweg7yaEnR36bj1tWPNYjVT6qwn-NHKOASQg0Kw4C8SrRMcyCMhvqmXkLbMWTaBVfYdwpKknHYtwegpnKG9Up1X1hV9uNEQnSQMVw5IrMvzSDhJV3-psgyB4frv5TKgXadOpIOHPBRykV3r4ZblV_zAOIkODA2nsgKe4Z-dsyTlVbgEV7Sn7GWBKa4l1Uf76WAzQ5xpnLauE5_oh4tT942xnA.5LcadxgSov-PQ6RNJoGeWg`,
			standCommonOutput)
		suite.Error(errParse, "token is expired")
	})
}

func (suite *JWEES256Suite) TestEES256JWTParseTokenMethodExpireNoExpiredTime() {
	suite.NotPanics(func() {
		j, err := NewEES256JWT(suite.key)
		suite.Equal(nil, err)

		standCommonOutput := NewClaimsBuilder().Build()
		errParse := j.VerifyToken(
			`eyJhbGciOiJFQ0RILUVTK0EyNTZLVyIsImN0eSI6IkpXVCIsImVuYyI6IkEyNTZHQ00iLCJlcGsiOnsia3R5IjoiRUMiLCJjcnYiOiJQLTI1NiIsIngiOiJRZXFMQmVVLVpMZXcxU2ZiaTJ6Wk9tb2dLdGdwMWJFZzFCVVFMVGNhVWtVIiwieSI6ImhoV2g2bEhPc3VQQlRIZWhxUlRhQXd6d2NEOGw1UVYydzdaX1RVamNXN0UifSwidHlwIjoiSldUIn0.Rb5kb8KYSYx2xWL_s66jTGkOjJr34BZDOIukameDwSVM-hdKrH_hPA.rxORxeUK2vC1Eur3.eoDMpC6i32nfQVcrmLE2dHCnjISjQtDyUtJmIAerBOzbHu0G53f8Z1iSMgNoAq45xTHM_xkU-pg013QVWjHbhGCCaGS1e18qslGg3GtiChtuDwmJcplLU4Q2G0vmiuXLicNxQAvzHZyO74TmJhLjT01Uo953zdSH90T5PlVYAN7zpLD6wd9ZUQ41KUCpeDUI6yBTe8yh3ReS9vJCVSLXJLQ06WvEp7vj2aCzDJDwB8-Ex8jDIwBPuxL_eTO9pFdQ4eyMj2kpACtU2IRtPPbLUJ5nOUN-aFzO966E07HvmAyOlWIpLsJPTZ8lSuUvcxpXCox9XZUWpJyHKACfhH4Vnfik9PuSeg.ozLrHUqfW8dn642WBFPq3g`,
			standCommonOutput)
		suite.NoError(errParse)
		suite.T().Log(standCommonOutput)
	})
}

func (suite *JWEES256Suite) TestEES256JWTRefreshTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewEES256JWT(suite.key)
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

func (suite *JWEES256Suite) TestEES256JWTRefreshTokenMethodExpire() {
	suite.NotPanics(func() {
		j, err := NewEES256JWT(suite.key)
		suite.Equal(nil, err)
		standClaimsOutput := NewClaimsBuilder().Build()
		tk, errRefresh := j.RefreshToken(
			"eyJhbGciOiJFQ0RILUVTK0EyNTZLVyIsImN0eSI6IkpXVCIsImVuYyI6IkEyNTZHQ00iLCJlcGsiOnsia3R5IjoiRUMiLCJjcnYiOiJQLTI1NiIsIngiOiJRZXFMQmVVLVpMZXcxU2ZiaTJ6Wk9tb2dLdGdwMWJFZzFCVVFMVGNhVWtVIiwieSI6ImhoV2g2bEhPc3VQQlRIZWhxUlRhQXd6d2NEOGw1UVYydzdaX1RVamNXN0UifSwidHlwIjoiSldUIn0.Rb5kb8KYSYx2xWL_s66jTGkOjJr34BZDOIukameDwSVM-hdKrH_hPA.rxORxeUK2vC1Eur3.eoDMpC6i32nfQVcrmLE2dHCnjISjQtDyUtJmIAerBOzbHu0G53f8Z1iSMgNoAq45xTHM_xkU-pg013QVWjHbhGCCaGS1e18qslGg3GtiChtuDwmJcplLU4Q2G0vmiuXLicNxQAvzHZyO74TmJhLjT01Uo953zdSH90T5PlVYAN7zpLD6wd9ZUQ41KUCpeDUI6yBTe8yh3ReS9vJCVSLXJLQ06WvEp7vj2aCzDJDwB8-Ex8jDIwBPuxL_eTO9pFdQ4eyMj2kpACtU2IRtPPbLUJ5nOUN-aFzO966E07HvmAyOlWIpLsJPTZ8lSuUvcxpXCox9XZUWpJyHKACfhH4Vnfik9PuSeg.ozLrHUqfW8dn642WBFPq3g",
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

func (suite *JWEES256Suite) TestEES256JWTRefreshTokenMethodparseRawError() {
	defer gostub.StubFunc(&parseSignedAndEncrypted, nil, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		j, err := NewEES256JWT(suite.key)
		suite.Equal(nil, err)
		standClaimsOutput := NewClaimsBuilder().Build()
		_, errRefresh := j.RefreshToken(
			"eyJhbGciOiJFQ0RILUVTK0EyNTZLVyIsImN0eSI6IkpXVCIsImVuYyI6IkEyNTZHQ00iLCJlcGsiOnsia3R5IjoiRUMiLCJjcnYiOiJQLTI1NiIsIngiOiJRZmJVRnNKM3NGV3ZKMkkxRGtQUElmYjQ3dGRTSmxsUzB3VjJ1VHhuQ2ZJIiwieSI6IjlXelNfTE5tZGkybzZUb3BVMXNwTnpfVlUwVk1QS3JZWkV4c0s5ODhRdHcifSwidHlwIjoiSldUIn0.AZwfbZgE2ib6ANat7lOVBBNfcE31Uccy_YlCedxkchL60xDxYAJqkA.lnzjMHoNKuNU57rZ.b6Vsp3g-ntbKH1J2XWkUrp9k7NCf-oNf6cTQ6lI800FAaplNIuRaZwQ_D7Rh3dUkSyNToYICgDKclNbVAivqTrqhbh6VJcUI7hbHMWIVCMu3Jh3TA6ZigwlrB9NZOZknvljpjLVMb7DZLZwC3rZnyinNHQTPISRbKO4Xqoe1g90_UwCEokG8UNm9GFsrudrfShunirj9xv42Z0aIKyDefrEsT3dQ3pGeX9uJl09S6RHYgEhLAYZ463_n7wixFf6OiPEuTIeQEp7nDi8HXGJkW3YIn0RbLK5V_RIvxCHgYX8SHlCYEzrAi97YjjplUdM._njneQ_kBKGFueqEjPdL5g",
			standClaimsOutput,
			100*time.Millisecond,
		)
		suite.Error(errRefresh, "got error")
	})
}

func TestJWEES256Suite(t *testing.T) {
	suite.Run(t, new(JWEES256Suite))
}
