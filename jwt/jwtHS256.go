package jwt

import (
	"github.com/cockroachdb/errors"
	"github.com/justdomepaul/toolbox/config"
	"github.com/square/go-jose/v3"
	"os"
	"time"
)

// NewHS256JWT method
func NewHS256JWT(hmacSecretKey string) (*HS256JWT, error) {
	sig, err := newSigner(
		jose.SigningKey{Algorithm: jose.HS256, Key: []byte(hmacSecretKey)},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		return nil, err
	}

	return &HS256JWT{
		SigningKey: []byte(hmacSecretKey),
		Sig:        sig,
	}, nil
}

func NewHS256JWTFromOptions(option config.JWT) (*HS256JWT, error) {
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
	return NewHS256JWT(key)
}

// HS256JWT type
type HS256JWT struct {
	SigningKey []byte
	Sig        jose.Signer
}

// GenerateToken method
func (h HS256JWT) GenerateToken(claims IJWTClaims) (token string, err error) {
	return signed(h.Sig).Claims(claims).CompactSerialize()
}

// Validate method
func (h HS256JWT) Validate(raw string) error {
	tok, errParse := parseSigned(raw)
	if errParse != nil {
		return errParse
	}
	return tok.Claims(h.SigningKey)
}

// VerifyToken method
func (h HS256JWT) VerifyToken(token string, claims IJWTClaims) (err error) {
	return parseHSRaw(h.SigningKey, token, claims)
}

// RefreshToken method
func (h HS256JWT) RefreshToken(token string, claims IJWTClaims, duration time.Duration) (string, error) {
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
