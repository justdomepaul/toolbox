package jwt

import (
	"crypto/rsa"
	"github.com/cockroachdb/errors"
	"github.com/justdomepaul/toolbox/config"
	"github.com/square/go-jose/v3"
	"os"
	"time"
)

// NewERS256JWT method
func NewERS256JWT(rsaPrivateKey string) (*ERS256JWT, error) {
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

	enc, err := newEncrypter(
		jose.A256GCM,
		jose.Recipient{Algorithm: jose.RSA_OAEP_256, Key: &rsaKey.PublicKey},
		(&jose.EncrypterOptions{}).WithType("JWT").WithContentType("JWT"),
	)
	if err != nil {
		return nil, err
	}

	return &ERS256JWT{
		SigningKey: rsaKey,
		Sig:        sig,
		Enc:        enc,
	}, nil
}

func NewERS256JWTFromOptions(option config.JWT) (*ERS256JWT, error) {
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
	return NewERS256JWT(key)
}

// ERS256JWT type
type ERS256JWT struct {
	SigningKey *rsa.PrivateKey
	Sig        jose.Signer
	Enc        jose.Encrypter
}

// GenerateToken method
func (er ERS256JWT) GenerateToken(claims IJWTClaims) (token string, err error) {
	return signedAndEncrypted(er.Sig, er.Enc).Claims(claims).CompactSerialize()
}

// Validate method
func (er ERS256JWT) Validate(raw string) error {
	tok, err := parseSignedAndEncrypted(raw)
	if err != nil {
		return err
	}
	_, errDecrypt := tok.Decrypt(er.SigningKey)
	if errDecrypt != nil {
		return errDecrypt
	}
	return nil
}

// VerifyToken method
func (er ERS256JWT) VerifyToken(token string, claims IJWTClaims) (err error) {
	return parseERSRaw(er.SigningKey, token, claims)
}

// RefreshToken method
func (er ERS256JWT) RefreshToken(token string, claims IJWTClaims, duration time.Duration) (string, error) {
	errParse := parseERSRaw(er.SigningKey, token, claims)
	if errParse != nil && errParse != ErrTokenExpired {
		return "", errParse
	}
	if errors.Is(errParse, ErrTokenExpired) {
		claims.(IJWTExpire).ExpiresAfter(duration)
		return signedAndEncrypted(er.Sig, er.Enc).Claims(claims).CompactSerialize()
	}
	return token, nil
}
