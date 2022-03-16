package jwt

import (
	"crypto/rsa"
	"github.com/cockroachdb/errors"
	"github.com/justdomepaul/toolbox/config"
	"github.com/square/go-jose/v3"
	"os"
	"time"
)

// NewRS256JWT method
func NewRS256JWT(rsaPrivateKey string) (*RS256JWT, error) {
	rsaKey, err := parseRSAPrivateKeyFromPEM([]byte(rsaPrivateKey))
	if err != nil {
		return nil, errors.Wrap(err, ErrParsePrivateKey.Error())
	}

	sig, err := newSigner(
		jose.SigningKey{Algorithm: jose.RS256, Key: rsaKey},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		return nil, err
	}

	return &RS256JWT{
		SigningKey: rsaKey,
		Sig:        sig,
	}, nil
}

func NewRS256JWTFromOptions(option config.JWT) (*RS256JWT, error) {
	if option.RsaPrivateKeyPath == "" && option.RsaPrivateKey == "" {
		return nil, ErrNoKey
	}
	key := option.RsaPrivateKey
	if option.RsaPrivateKeyPath != "" {
		if result, err := os.ReadFile(option.RsaPrivateKeyPath); err != nil {
			return nil, err
		} else {
			key = string(result)
		}
	}
	return NewRS256JWT(key)
}

// RS256JWT type
type RS256JWT struct {
	SigningKey *rsa.PrivateKey
	Sig        jose.Signer
}

// GenerateToken method
func (r RS256JWT) GenerateToken(claims IJWTClaims) (token string, err error) {
	return signed(r.Sig).Claims(claims).CompactSerialize()
}

// Validate method
func (r RS256JWT) Validate(raw string) error {
	tok, errParse := parseSigned(raw)
	if errParse != nil {
		return errParse
	}
	return tok.Claims(r.SigningKey.Public())
}

// VerifyToken method
func (r RS256JWT) VerifyToken(token string, claims IJWTClaims) (err error) {
	return parseRSRaw(r.SigningKey, token, claims)
}

// RefreshToken method
func (r RS256JWT) RefreshToken(token string, claims IJWTClaims, duration time.Duration) (string, error) {
	errParse := parseRSRaw(r.SigningKey, token, claims)
	if errParse != nil && errParse != ErrTokenExpired {
		return "", errParse
	}
	if errors.Is(errParse, ErrTokenExpired) {
		claims.(IJWTExpire).ExpiresAfter(duration)
		return signed(r.Sig).Claims(claims).CompactSerialize()
	}
	return token, nil
}
