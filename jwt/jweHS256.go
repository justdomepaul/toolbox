package jwt

import (
	"github.com/cockroachdb/errors"
	"github.com/justdomepaul/toolbox/config"
	"github.com/square/go-jose/v3"
	"os"
	"time"
)

// NewEHS256JWT method
func NewEHS256JWT(hmacSecretKey string) (*EHS256JWT, error) {
	sig, err := newSigner(
		jose.SigningKey{Algorithm: jose.HS256, Key: []byte(hmacSecretKey)},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		return nil, err
	}

	enc, err := newEncrypter(
		jose.A128CBC_HS256,
		jose.Recipient{Algorithm: jose.PBES2_HS256_A128KW, Key: []byte(hmacSecretKey)},
		(&jose.EncrypterOptions{}).WithType("JWT").WithContentType("JWT"),
	)
	if err != nil {
		return nil, err
	}

	return &EHS256JWT{
		SigningKey: []byte(hmacSecretKey),
		Sig:        sig,
		Enc:        enc,
	}, nil
}

func NewEHS256JWTFromOptions(option config.JWT) (*EHS256JWT, error) {
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
	return NewEHS256JWT(key)
}

// EHS256JWT type
type EHS256JWT struct {
	SigningKey []byte
	Sig        jose.Signer
	Enc        jose.Encrypter
}

// GenerateToken method
func (eh EHS256JWT) GenerateToken(claims IJWTClaims) (token string, err error) {
	return signedAndEncrypted(eh.Sig, eh.Enc).Claims(claims).CompactSerialize()
}

// Validate method
func (eh EHS256JWT) Validate(raw string) error {
	tok, err := parseSignedAndEncrypted(raw)
	if err != nil {
		return err
	}
	_, errDecrypt := tok.Decrypt(eh.SigningKey)
	if errDecrypt != nil {
		return errDecrypt
	}
	return nil
}

// VerifyToken method
func (eh EHS256JWT) VerifyToken(token string, claims IJWTClaims) (err error) {
	return parseEHSRaw(eh.SigningKey, token, claims)
}

// RefreshToken method
func (eh EHS256JWT) RefreshToken(token string, claims IJWTClaims, duration time.Duration) (string, error) {
	errParse := parseEHSRaw(eh.SigningKey, token, claims)
	if errParse != nil && errParse != ErrTokenExpired {
		return "", errParse
	}
	if errors.Is(errParse, ErrTokenExpired) {
		claims.(IJWTExpire).ExpiresAfter(duration)
		return signedAndEncrypted(eh.Sig, eh.Enc).Claims(claims).CompactSerialize()
	}
	return token, nil
}
