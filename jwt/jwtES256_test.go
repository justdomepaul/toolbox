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

type JWTES256Suite struct {
	suite.Suite
	key    string
	option config.JWT
}

func (suite *JWTES256Suite) SetupTest() {
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

func (suite *JWTES256Suite) TestNewES256JWT() {
	suite.NotPanics(func() {
		_, err := NewES256JWT(suite.key)
		suite.Equal(nil, err)
	})
}

func (suite *JWTES256Suite) TestNewES256JWTFromOptions() {
	suite.NotPanics(func() {
		_, err := NewES256JWTFromOptions(suite.option)
		suite.NoError(err)
	})
}

func (suite *JWTES256Suite) TestNewES256JWTFromOptionsAllNoOption() {
	suite.NotPanics(func() {
		_, err := NewES256JWTFromOptions(config.JWT{
			EcdsaPrivateKeyPath: "",
			EcdsaPrivateKey:     "",
		})
		suite.Error(err, ErrNoKey)
	})
}

func (suite *JWTES256Suite) TestNewES256JWTFromOptionsAllNoOptionNoFIle() {
	suite.NotPanics(func() {
		_, err := NewES256JWTFromOptions(config.JWT{
			EcdsaPrivateKeyPath: "testFile.txt",
			EcdsaPrivateKey:     "",
		})
		suite.Error(err)
	})
}

func (suite *JWTES256Suite) TestNewES256JWTparseECPrivateKeyFromPEMError() {
	defer gostub.StubFunc(&parseECPrivateKeyFromPEM, nil, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewES256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWTES256Suite) TestNewES256JWTnewSignerError() {
	tSigner := testSigner{}
	defer gostub.StubFunc(&newSigner, tSigner, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		_, err := NewES256JWT(suite.key)
		suite.Error(err, "got error")
	})
}

func (suite *JWTES256Suite) TestES256JWTGenerateTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewES256JWT(suite.key)
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

func (suite *JWTES256Suite) TestES256JWTValidateMethod() {
	suite.NotPanics(func() {
		j, err := NewES256JWT(suite.key)
		suite.Equal(nil, err)

		suite.NoError(j.Validate(
			"eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0ZXJDbGllbnQiLCJleHAiOjE2NDY5NzE5NzQsImlzcyI6InRlc3RlciIsImp0aSI6InRlc3QwMDEiLCJzIjoidGVzdERhdGEiLCJzdWIiOiJ0ZXN0VG9waWMifQ.Dpf_yj1NUWfLOpo5isJkSbkvhj5jBs-mbx7NokyzRDlWPxMHT3Dmi6dO5QqhfTFpogEsP0cde__iD5Yyi0mxEA"))
	})
}

func (suite *JWTES256Suite) TestES256JWTValidateMethodparseSignedError() {
	suite.NotPanics(func() {
		defer gostub.StubFunc(&parseSigned, nil, errors.New("got error")).Reset()
		j, err := NewES256JWT(suite.key)
		suite.Equal(nil, err)

		suite.Error(j.Validate(
			"eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0ZXJDbGllbnQiLCJleHAiOjE2NDY5NzE5NzQsImlzcyI6InRlc3RlciIsImp0aSI6InRlc3QwMDEiLCJzIjoidGVzdERhdGEiLCJzdWIiOiJ0ZXN0VG9waWMifQ.Dpf_yj1NUWfLOpo5isJkSbkvhj5jBs-mbx7NokyzRDlWPxMHT3Dmi6dO5QqhfTFpogEsP0cde__iD5Yyi0mxEA"))
	})
}

func (suite *JWTES256Suite) TestES256JWTParseTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewES256JWT(suite.key)
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

func (suite *JWTES256Suite) TestES256JWTParseTokenMethodExpire() {
	suite.NotPanics(func() {
		j, err := NewES256JWT(suite.key)
		suite.Equal(nil, err)

		standCommonOutput := NewCommon(NewClaimsBuilder().Build())
		errParse := j.VerifyToken(
			"eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0ZXJDbGllbnQiLCJleHAiOjE2MzMzMjAxMzYsImlzcyI6InRlc3RlciIsImp0aSI6InRlc3QwMDEiLCJzIjoidGVzdERhdGEiLCJzdWIiOiJ0ZXN0VG9waWMifQ.h5E24wuAu4BNLajw8I_PO2VYNHht5SvNf0EuvK5w1WQPAa4LDmfKYHYgEzyuZOTM65B1IdjQvvCCDbRwG4YAOA",
			standCommonOutput)
		suite.Error(errParse, "token is expired")
	})
}

func (suite *JWTES256Suite) TestES256JWTParseTokenMethodExpireNoExpired() {
	suite.NotPanics(func() {
		j, err := NewES256JWT(suite.key)
		suite.Equal(nil, err)

		standCommonOutput := NewClaimsBuilder().Build()
		errParse := j.VerifyToken(
			"eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0ZXJDbGllbnQiLCJpc3MiOiJ0ZXN0ZXIiLCJqdGkiOiJ0ZXN0MDAxIiwicyI6InRlc3REYXRhIiwic3ViIjoidGVzdFRvcGljIn0.XI3qtSbFJKieZgPK8Z4x9DhEBYHq20zhtwVv6pdI6EQtECDTRMiaioaX077F6nlmhiIf6i_cK920AqEKFMlrpw",
			standCommonOutput)
		suite.NoError(errParse)
	})
}

func (suite *JWTES256Suite) TestES256JWTRefreshTokenMethod() {
	suite.NotPanics(func() {
		j, err := NewES256JWT(suite.key)
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

func (suite *JWTES256Suite) TestES256JWTRefreshTokenMethodExpire() {
	suite.NotPanics(func() {
		j, err := NewES256JWT(suite.key)
		suite.Equal(nil, err)
		standClaimsOutput := NewClaimsBuilder().Build()
		tk, errRefresh := j.RefreshToken(
			"eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0ZXJDbGllbnQiLCJleHAiOjE2NDY5NzIwMzksImlzcyI6InRlc3RlciIsImp0aSI6InRlc3QwMDEiLCJzIjoidGVzdERhdGEiLCJzdWIiOiJ0ZXN0VG9waWMifQ.u6Y3aAgY0Qv9Y5jHrYtKn-R2F6y3DygPlndJCvSyXMs7f8cP1rLFJZAKM07rIt9zIeZ-Mo2RdJNtzQK-kxWh6g",
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

func (suite *JWTES256Suite) TestES256JWTRefreshTokenMethodparseRawError() {
	defer gostub.StubFunc(&parseSigned, nil, errors.New("got error")).Reset()
	suite.NotPanics(func() {
		j, err := NewES256JWT(suite.key)
		suite.Equal(nil, err)
		standClaimsOutput := NewClaimsBuilder().Build()
		_, errRefresh := j.RefreshToken(
			"eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ0ZXN0ZXJDbGllbnQiLCJleHAiOjE2NDY5NzE5NzQsImlzcyI6InRlc3RlciIsImp0aSI6InRlc3QwMDEiLCJzIjoidGVzdERhdGEiLCJzdWIiOiJ0ZXN0VG9waWMifQ.Dpf_yj1NUWfLOpo5isJkSbkvhj5jBs-mbx7NokyzRDlWPxMHT3Dmi6dO5QqhfTFpogEsP0cde__iD5Yyi0mxEA",
			standClaimsOutput,
			100*time.Millisecond,
		)
		suite.Error(errRefresh, "got error")
	})
}

func TestJWTES256Suite(t *testing.T) {
	suite.Run(t, new(JWTES256Suite))
}
