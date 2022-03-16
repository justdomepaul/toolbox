package jwt

import (
	"crypto/ecdsa"
	"github.com/cockroachdb/errors"
	"github.com/justdomepaul/toolbox/config"
	"github.com/square/go-jose/v3"
	"os"
	"time"
)

// NewEES256JWT method
func NewEES256JWT(ecdsaPrivateKey string) (*EES256JWT, error) {
	ecdsaKey, err := parseECPrivateKeyFromPEM([]byte(ecdsaPrivateKey))
	if err != nil {
		return nil, errors.Wrap(err, ErrParsePrivateKey.Error())
	}

	sig, err := newSigner(
		jose.SigningKey{Algorithm: jose.ES256, Key: ecdsaKey},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		return nil, err
	}

	enc, err := newEncrypter(
		jose.A256GCM,
		jose.Recipient{Algorithm: jose.ECDH_ES_A256KW, Key: &ecdsaKey.PublicKey},
		(&jose.EncrypterOptions{}).WithType("JWT").WithContentType("JWT"),
	)
	if err != nil {
		return nil, err
	}

	return &EES256JWT{
		SigningKey: ecdsaKey,
		Sig:        sig,
		Enc:        enc,
	}, nil
}

func NewEES256JWTFromOptions(option config.JWT) (*EES256JWT, error) {
	if option.EcdsaPrivateKeyPath == "" && option.EcdsaPrivateKey == "" {
		return nil, ErrNoKey
	}
	key := option.EcdsaPrivateKey
	if option.EcdsaPrivateKeyPath != "" {
		if result, err := os.ReadFile(option.EcdsaPrivateKeyPath); err != nil {
			return nil, err
		} else {
			key = string(result)
		}
	}
	return NewEES256JWT(key)
}

// EES256JWT type
type EES256JWT struct {
	SigningKey *ecdsa.PrivateKey
	Sig        jose.Signer
	Enc        jose.Encrypter
}

// GenerateToken method
func (ee EES256JWT) GenerateToken(claims IJWTClaims) (token string, err error) {
	return signedAndEncrypted(ee.Sig, ee.Enc).Claims(claims).CompactSerialize()
}

// Validate method
func (ee EES256JWT) Validate(raw string) error {
	tok, err := parseSignedAndEncrypted(raw)
	if err != nil {
		return err
	}
	_, errDecrypt := tok.Decrypt(ee.SigningKey)
	if errDecrypt != nil {
		return errDecrypt
	}
	return nil
}

// VerifyToken method
func (ee EES256JWT) VerifyToken(token string, claims IJWTClaims) (err error) {
	return parseEESRaw(ee.SigningKey, token, claims)
}

// RefreshToken method
func (ee EES256JWT) RefreshToken(token string, claims IJWTClaims, duration time.Duration) (string, error) {
	errParse := parseEESRaw(ee.SigningKey, token, claims)
	if errParse != nil && errParse != ErrTokenExpired {
		return "", errParse
	}
	if errors.Is(errParse, ErrTokenExpired) {
		claims.(IJWTExpire).ExpiresAfter(duration)
		return signedAndEncrypted(ee.Sig, ee.Enc).Claims(claims).CompactSerialize()
	}
	return token, nil
}
