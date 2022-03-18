package jwt

import (
	"github.com/cockroachdb/errors"
	"github.com/go-jose/go-jose/v3"
	"github.com/justdomepaul/toolbox/config"
	"os"
	"time"
)

// NewHS384JWT method
func NewHS384JWT(hmacSecretKey string) (*HS384JWT, error) {
	sig, err := newSigner(
		jose.SigningKey{Algorithm: jose.HS384, Key: []byte(hmacSecretKey)},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		return nil, err
	}

	return &HS384JWT{
		SigningKey: []byte(hmacSecretKey),
		Sig:        sig,
	}, nil
}

func NewHS384JWTFromOptions(option config.JWT) (*HS384JWT, error) {
	if option.HmacSecretKeyPath == "" && option.HmacSecretKey == "" {
		return nil, ErrNoKey
	}
	key := option.HmacSecretKey
	if option.HmacSecretKeyPath != "" {
		if result, err := os.ReadFile(option.HmacSecretKeyPath); err != nil {
			return nil, err
		} else {
			key = string(result)
		}
	}
	return NewHS384JWT(key)
}

// HS384JWT type
type HS384JWT struct {
	SigningKey []byte
	Sig        jose.Signer
}

// GenerateToken method
func (h HS384JWT) GenerateToken(claims IJWTClaims) (token string, err error) {
	return signed(h.Sig).Claims(claims).CompactSerialize()
}

// Validate method
func (h HS384JWT) Validate(raw string) error {
	tok, errParse := parseSigned(raw)
	if errParse != nil {
		return errParse
	}
	return tok.Claims(h.SigningKey)
}

// VerifyToken method
func (h HS384JWT) VerifyToken(token string, claims IJWTClaims) (err error) {
	return parseHSRaw(h.SigningKey, token, claims)
}

// RefreshToken method
func (h HS384JWT) RefreshToken(token string, claims IJWTClaims, duration time.Duration) (string, error) {
	errParse := parseHSRaw(h.SigningKey, token, claims)
	if errParse != nil && errParse != ErrTokenExpired {
		return "", errParse
	}
	if errors.Is(errParse, ErrTokenExpired) {
		claims.(IJWTExpire).ExpiresAfter(duration)
		return signed(h.Sig).Claims(claims).CompactSerialize()
	}
	return token, nil
}
