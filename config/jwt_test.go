package config

import (
	"github.com/stretchr/testify/assert"
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
	t := suite.T()
	os.Clearenv()
	suite.EcdsaPrivateKeyPath = "testEcdsaPrivateKeyPath"
	suite.EcdsaPrivateKey = "testEcdsaPrivateKey"
	suite.RsaPrivateKeyPath = "testRsaPrivateKeyPath"
	suite.RsaPrivateKey = "testRsaPrivateKey"
	suite.HmacSecretKeyPath = "testHmacSecretKeyPath"
	suite.HmacSecretKey = "testHmacSecretKey"

	assert.NoError(t, os.Setenv("ECDSA_PRIVATE_KEY_PATH", suite.EcdsaPrivateKeyPath))
	assert.NoError(t, os.Setenv("ECDSA_PRIVATE_KEY", suite.EcdsaPrivateKey))
	assert.NoError(t, os.Setenv("RSA_PRIVATE_KEY_PATH", suite.RsaPrivateKeyPath))
	assert.NoError(t, os.Setenv("RSA_PRIVATE_KEY", suite.RsaPrivateKey))
	assert.NoError(t, os.Setenv("HMAC_SECRET_KEY_PATH", suite.HmacSecretKeyPath))
	assert.NoError(t, os.Setenv("HMAC_SECRET_KEY", suite.HmacSecretKey))
}

func (suite *JWTSuite) TestDefaultOption() {
	t := suite.T()
	jwt := &JWT{}
	suite.NoError(LoadFromEnv(jwt))
	assert.Equal(t, suite.EcdsaPrivateKeyPath, jwt.EcdsaPrivateKeyPath)
	assert.Equal(t, suite.EcdsaPrivateKey, jwt.EcdsaPrivateKey)
	assert.Equal(t, suite.RsaPrivateKeyPath, jwt.RsaPrivateKeyPath)
	assert.Equal(t, suite.RsaPrivateKey, jwt.RsaPrivateKey)
	assert.Equal(t, suite.HmacSecretKeyPath, jwt.HmacSecretKeyPath)
	assert.Equal(t, suite.HmacSecretKey, jwt.HmacSecretKey)
}

func TestJWTSuite(t *testing.T) {
	suite.Run(t, new(JWTSuite))
}
