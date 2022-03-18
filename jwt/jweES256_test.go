package jwt

import (
	"github.com/cockroachdb/errors"
	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/justdomepaul/toolbox/config"
	"github.com/prashantv/gostub"
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
	jwt    IJWT
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

	j, err := NewEES256JWT(suite.key)
	suite.NoError(err)
	suite.jwt = j
}

func (suite *JWEES256Suite) TestJWT() {
	suite.NotPanics(func() {
		_, err := NewEES256JWT(suite.key)
		suite.NoError(err)
	})
}

func (suite *JWEES256Suite) TestJWTFromOptions() {
	suite.NotPanics(func() {
		_, err := NewEES256JWTFromOptions(suite.option)
		suite.NoError(err)
	})
}

func (suite *JWEES256Suite) TestJWTFromOptionsAllNoOption() {
	suite.NotPanics(func() {
		_, err := NewEES256JWTFromOptions(config.JWT{
			EcdsaPrivateKeyPath: "",
			EcdsaPrivateKey:     "",
		})
		suite.Error(err, ErrNoKey)
	})
}

func (suite *JWEES256Suite) TestJWTFromOptionsAllNoOptionNoFIle() {
	suite.NotPanics(func() {
		_, err := NewEES256JWTFromOptions(config.JWT{
			EcdsaPrivateKeyPath: "testFile.txt",
			EcdsaPrivateKey:     "",
		})
		suite.Error(err)
	})
}

func (suite *JWEES256Suite) TestJWTParseECPrivateKeyFromPEMError() {
	defer gostub.StubFunc(&parseECPrivateKeyFromPEM, nil, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewEES256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWEES256Suite) TestJWTNewSignerError() {
	tSigner := testSigner{}
	defer gostub.StubFunc(&newSigner, tSigner, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewEES256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWEES256Suite) TestJWTNewEncrypterError() {
	tEncrypter := testEncrypter{}
	defer gostub.StubFunc(&newEncrypter, tEncrypter, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewEES256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWEES256Suite) TestJWTGenerateTokenMethod() {
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

func (suite *JWEES256Suite) TestJWTValidateMethod() {
	suite.NotPanics(func() {
		suite.NoError(suite.jwt.Validate(
			"eyJhbGciOiJFQ0RILUVTK0EyNTZLVyIsImN0eSI6IkpXVCIsImVuYyI6IkEyNTZHQ00iLCJlcGsiOnsia3R5IjoiRUMiLCJjcnYiOiJQLTI1NiIsIngiOiI0QTBYYjVSbFgtV0hRUm9wcDdkdkp2THgxTEpxMTl6TFdvcVE3NDQzTndnIiwieSI6ImNqYjZBckJDZ0QtZUl0ZnJ2b2RUdmtjVEdQMlhGNU9nVkJyTHVZcUlfWW8ifSwidHlwIjoiSldUIn0.AxPUTCR7fYQlUU093IO3p5v1lnATiazofKI_A3AOW8TqHvTknehLrw.Qy_ZqSBj8CGPUw-W.B-MCpyCLcRBHX8lOjhk2Wv-T5lsj1zUUIiqHdQAkiVyff8dFoChpP5MUh2YakP9ZU6FKFxQdgEuHMNTj9p8wOSn1Oeq_BJ3xQcr4sDq26TlXXRIj-30h05BpzqPrBNWKuDIAx0lkPjCfVgAs4YMx7e9hX5B5Fs28LFb76nlRQUITu3TiPwnxJqaMyILu8JqOsxnJ5XnMbypA_7xai5mG7oMhsHfvQCVnLm0Mly3vGLILtOsdlXVOc7OTbO0rpJ_ifxAr9l0YTIlMcZe1X9m3ER6FQSejs-vJ5jvwwkIl61QJEjxjvKWfR2vXW_XWcc4AWaivJaMpzFwkiZ5QgSpVycWcFVitVA.9vNJQObVc0v8f-YS3g5TcA"))
	})
}

func (suite *JWEES256Suite) TestJWTValidateMethodParseSignedAndEncryptedError() {
	suite.NotPanics(func() {
		defer gostub.StubFunc(&parseSignedAndEncrypted, nil, errors.New("got error")).Reset()

		suite.Error(suite.jwt.Validate(""))
	})
}

func (suite *JWEES256Suite) TestJWTVerifyTokenMethod() {
	suite.NotPanics(func() {
		standClaims := NewClaimsBuilder().
			WithSubject("testTopic").
			WithIssuer("tester").
			WithID("test001").
			WithAudience([]string{"testerClient"}).
			ExpiresAfter(3 * time.Second)
		common := NewCommon(standClaims.Build(), WithSecret("testData"))
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

func (suite *JWEES256Suite) TestJWTVerifyTokenMethodExpire() {
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

func (suite *JWEES256Suite) TestJWTVerifyTokenMethodExpireNoExpiredTime() {
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

func (suite *JWEES256Suite) TestJWTRefreshTokenMethod() {
	suite.NotPanics(func() {
		standClaims := NewClaimsBuilder().
			WithSubject("testTopic").
			WithIssuer("tester").
			WithID("test001").
			WithAudience([]string{"testerClient"}).
			ExpiresAfter(3 * time.Second)
		common := NewCommon(standClaims.Build(), WithSecret("testData"))
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

func (suite *JWEES256Suite) TestJWTRefreshTokenMethodExpire() {
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

func (suite *JWEES256Suite) TestJWTRefreshTokenMethodExpireNoIJWTExpireClaim() {
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

func (suite *JWEES256Suite) TestJWTRefreshTokenMethodParseRawError() {
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

func TestJWEES256Suite(t *testing.T) {
	suite.Run(t, new(JWEES256Suite))
}
