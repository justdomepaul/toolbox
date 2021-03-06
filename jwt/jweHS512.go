package jwt

import (
	"github.com/cockroachdb/errors"
	"github.com/go-jose/go-jose/v3"
	"github.com/justdomepaul/toolbox/config"
	"os"
	"time"
)

// NewEHS512JWT method
func NewEHS512JWT(hmacSecretKey string) (*EHS512JWT, error) {
	sig, err := newSigner(
		jose.SigningKey{Algorithm: jose.HS512, Key: []byte(hmacSecretKey)},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		return nil, err
	}

	enc, err := newEncrypter(
		jose.A256CBC_HS512,
		jose.Recipient{Algorithm: jose.PBES2_HS512_A256KW, Key: []byte(hmacSecretKey)},
		(&jose.EncrypterOptions{}).WithType("JWT").WithContentType("JWT"),
	)
	if err != nil {
		return nil, err
	}

	return &EHS512JWT{
		SigningKey: []byte(hmacSecretKey),
		Sig:        sig,
		Enc:        enc,
	}, nil
}

func NewEHS512JWTFromOptions(option config.JWT) (*EHS512JWT, error) {
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
	return NewEHS512JWT(key)
}

// EHS512JWT type
type EHS512JWT struct {
	SigningKey []byte
	Sig        jose.Signer
	Enc        jose.Encrypter
}

// GenerateToken method
func (eh EHS512JWT) GenerateToken(claims IJWTClaims) (token string, err error) {
	return signedAndEncrypted(eh.Sig, eh.Enc).Claims(claims).CompactSerialize()
}

// Validate method
func (eh EHS512JWT) Validate(raw string) error {
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
func (eh EHS512JWT) VerifyToken(token string, claims IJWTClaims) (err error) {
	return parseEHSRaw(eh.SigningKey, token, claims)
}

// RefreshToken method
func (eh EHS512JWT) RefreshToken(token string, claims IJWTClaims, duration time.Duration) (string, error) {
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
