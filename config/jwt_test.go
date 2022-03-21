package config

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type JWTSuite struct {
	suite.Suite
	EcdsaPrivateKeyPath string
	EcdsaPrivateKey     string
	RsaPrivateKeyPath   string
	RsaPrivateKey       string
	HmacSecretKeyPath   string
	HmacSecretKey       string
}

func (suite *JWTSuite) SetupSuite() {
	os.Clearenv()
	suite.EcdsaPrivateKeyPath = "testEcdsaPrivateKeyPath"
	suite.EcdsaPrivateKey = "testEcdsaPrivateKey"
	suite.RsaPrivateKeyPath = "testRsaPrivateKeyPath"
	suite.RsaPrivateKey = "testRsaPrivateKey"
	suite.HmacSecretKeyPath = "testHmacSecretKeyPath"
	suite.HmacSecretKey = "testHmacSecretKey"

	suite.NoError(os.Setenv("ECDSA_PRIVATE_KEY_PATH", suite.EcdsaPrivateKeyPath))
	suite.NoError(os.Setenv("ECDSA_PRIVATE_KEY", suite.EcdsaPrivateKey))
	suite.NoError(os.Setenv("RSA_PRIVATE_KEY_PATH", suite.RsaPrivateKeyPath))
	suite.NoError(os.Setenv("RSA_PRIVATE_KEY", suite.RsaPrivateKey))
	suite.NoError(os.Setenv("HMAC_SECRET_KEY_PATH", suite.HmacSecretKeyPath))
	suite.NoError(os.Setenv("HMAC_SECRET_KEY", suite.HmacSecretKey))
}

func (suite *JWTSuite) TestDefaultOption() {
	jwt := &JWT{}
	suite.NoError(LoadFromEnv(jwt))
	suite.Equal(suite.EcdsaPrivateKeyPath, jwt.EcdsaPrivateKeyPath)
	suite.Equal(suite.EcdsaPrivateKey, jwt.EcdsaPrivateKey)
	suite.Equal(suite.RsaPrivateKeyPath, jwt.RsaPrivateKeyPath)
	suite.Equal(suite.RsaPrivateKey, jwt.RsaPrivateKey)
	suite.Equal(suite.HmacSecretKeyPath, jwt.HmacSecretKeyPath)
	suite.Equal(suite.HmacSecretKey, jwt.HmacSecretKey)
}

func TestJWTSuite(t *testing.T) {
	suite.Run(t, new(JWTSuite))
}
