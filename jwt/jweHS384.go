package jwt

import (
	"github.com/justdomepaul/toolbox/config"
	"github.com/square/go-jose/v3"
	"os"
	"time"
)

// NewEHS384JWT method
func NewEHS384JWT(hmacSecretKey string) (*EHS384JWT, error) {
	sig, err := newSigner(
		jose.SigningKey{Algorithm: jose.HS384, Key: []byte(hmacSecretKey)},
		(&jose.SignerOptions{}).WithType("JWT"),
	)
	if err != nil {
		return nil, err
	}

	enc, err := newEncrypter(
		jose.A192CBC_HS384,
		jose.Recipient{Algorithm: jose.PBES2_HS384_A192KW, Key: []byte(hmacSecretKey)},
		(&jose.EncrypterOptions{}).WithType("JWT").WithContentType("JWT"),
	)
	if err != nil {
		return nil, err
	}

	return &EHS384JWT{
		SigningKey: []byte(hmacSecretKey),
		Sig:        sig,
		Enc:        enc,
	}, nil
}

func NewEHS384JWTFromOptions(option config.JWT) (*EHS384JWT, error) {
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
	return NewEHS384JWT(key)
}

// EHS384JWT type
type EHS384JWT struct {
	SigningKey []byte
	Sig        jose.Signer
	Enc        jose.Encrypter
}

// GenerateToken method
func (eh EHS384JWT) GenerateToken(claims IJWTClaims) (token string, err error) {
	return signedAndEncrypted(eh.Sig, eh.Enc).Claims(claims).CompactSerialize()
}

// Validate method
func (eh EHS384JWT) Validate(raw string) error {
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

// ParseToken method
func (eh EHS384JWT) VerifyToken(token string, claims IJWTClaims) (err error) {
	return parseEHSRaw(eh.SigningKey, token, claims)
}

// RefreshToken method
func (eh EHS384JWT) RefreshToken(token string, claims IJWTClaims, duration time.Duration) (string, error) {
	errParse := parseEHSRaw(eh.SigningKey, token, claims)
	if errParse != nil && errParse != ErrTokenExpired {
		return "", errParse
	}
	if errParse == ErrTokenExpired {
		claims.(IJWTExpire).ExpiresAfter(duration)
		return signedAndEncrypted(eh.Sig, eh.Enc).Claims(claims).CompactSerialize()
	}
	return token, nil
}
