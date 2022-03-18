package jwt

import (
	"crypto/ecdsa"
	"github.com/cockroachdb/errors"
	"github.com/go-jose/go-jose/v3"
	"github.com/justdomepaul/toolbox/config"
	"os"
	"time"
)

// NewES256JWT method
func NewES256JWT(ecdsaPrivateKey string) (*ES256JWT, error) {
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

	return &ES256JWT{
		SigningKey: ecdsaKey,
		Sig:        sig,
	}, nil
}

func NewES256JWTFromOptions(option config.JWT) (*ES256JWT, error) {
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
	return NewES256JWT(key)
}

// ES256JWT type
type ES256JWT struct {
	SigningKey *ecdsa.PrivateKey
	Sig        jose.Signer
}

// GenerateToken method
func (e ES256JWT) GenerateToken(claims IJWTClaims) (token string, err error) {
	return signed(e.Sig).Claims(claims).CompactSerialize()
}

// Validate method
func (e ES256JWT) Validate(raw string) error {
	tok, errParse := parseSigned(raw)
	if errParse != nil {
		return errParse
	}
	return tok.Claims(e.SigningKey.Public())
}

// VerifyToken method
func (e ES256JWT) VerifyToken(token string, claims IJWTClaims) (err error) {
	return parseESRaw(e.SigningKey, token, claims)
}

// RefreshToken method
func (e ES256JWT) RefreshToken(token string, claims IJWTClaims, duration time.Duration) (string, error) {
	errParse := parseESRaw(e.SigningKey, token, claims)
	if errParse != nil && errParse != ErrTokenExpired {
		return "", errParse
	}
	if errors.Is(errParse, ErrTokenExpired) {
		if instance, ok := claims.(IJWTExpire); ok {
			instance.ExpiresAfter(duration)
		}
		return signed(e.Sig).Claims(claims).CompactSerialize()
	}
	return token, nil
}
