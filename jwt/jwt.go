package jwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"github.com/cockroachdb/errors"
	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	jwtPkg "github.com/golang-jwt/jwt"
	"time"
)

var (
	ErrNoKey = errors.New("not provider jwt private/secret key")
	// ErrTokenExpired variable
	ErrTokenExpired = errors.New("token is expired")
	// ErrParsePrivateKey variable
	ErrParsePrivateKey = errors.New("parse private key error")
)

var (
	now                       = time.Now
	newSigner                 = jose.NewSigner
	newEncrypter              = jose.NewEncrypter
	parseECPrivateKeyFromPEM  = jwtPkg.ParseECPrivateKeyFromPEM
	signedAndEncrypted        = jwt.SignedAndEncrypted
	parseSignedAndEncrypted   = jwt.ParseSignedAndEncrypted
	parseRSAPrivateKeyFromPEM = jwtPkg.ParseRSAPrivateKeyFromPEM
	signed                    = jwt.Signed
	parseSigned               = jwt.ParseSigned
)

// IJWT interface
type IJWT interface {
	GenerateToken(claims IJWTClaims) (string, error)
	Validate(raw string) (err error)
	VerifyToken(token string, claims IJWTClaims) (err error)
	RefreshToken(token string, claims IJWTClaims, duration time.Duration) (string, error)
}

func ParseUnverified(raw string, claims IJWTClaims) error {
	tok, errParse := parseSigned(raw)
	if errParse != nil {
		return errParse
	}
	return tok.UnsafeClaimsWithoutVerification(claims)
}

func parseEESRaw(signingKey *ecdsa.PrivateKey, raw string, claims IJWTClaims) error {
	tok, errParse := parseSignedAndEncrypted(raw)
	if errParse != nil {
		return errParse
	}
	nested, errDecrypt := tok.Decrypt(signingKey)
	if errDecrypt != nil {
		return errDecrypt
	}
	errClaims := nested.Claims(signingKey.Public(), claims)
	if errClaims != nil {
		return errClaims
	}
	return checkExpire(claims)
}

func parseEHSRaw(signingKey []byte, raw string, claims IJWTClaims) error {
	tok, errParse := parseSignedAndEncrypted(raw)
	if errParse != nil {
		return errParse
	}
	nested, errDecrypt := tok.Decrypt(signingKey)
	if errDecrypt != nil {
		return errDecrypt
	}
	errClaims := nested.Claims(signingKey, claims)
	if errClaims != nil {
		return errClaims
	}
	return checkExpire(claims)
}

func parseERSRaw(signingKey *rsa.PrivateKey, raw string, claims IJWTClaims) error {
	tok, errParse := parseSignedAndEncrypted(raw)
	if errParse != nil {
		return errParse
	}
	nested, errDecrypt := tok.Decrypt(signingKey)
	if errDecrypt != nil {
		return errDecrypt
	}
	errClaims := nested.Claims(signingKey.Public(), claims)
	if errClaims != nil {
		return errClaims
	}
	return checkExpire(claims)
}

func parseESRaw(signingKey *ecdsa.PrivateKey, raw string, claims IJWTClaims) error {
	tok, errParse := parseSigned(raw)
	if errParse != nil {
		return errParse
	}

	errClaims := tok.Claims(signingKey.Public(), claims)
	if errClaims != nil {
		return errClaims
	}
	return checkExpire(claims)
}

func parseHSRaw(signingKey []byte, raw string, claims IJWTClaims) error {
	tok, errParse := parseSigned(raw)
	if errParse != nil {
		return errParse
	}

	errClaims := tok.Claims(signingKey, claims)
	if errClaims != nil {
		return errClaims
	}
	return checkExpire(claims)
}

func parseRSRaw(signingKey *rsa.PrivateKey, raw string, claims IJWTClaims) error {
	tok, errParse := parseSigned(raw)
	if errParse != nil {
		return errParse
	}

	errClaims := tok.Claims(signingKey.Public(), claims)
	if errClaims != nil {
		return errClaims
	}
	return checkExpire(claims)
}

func checkExpire(claims IJWTClaims) error {
	if instance, ok := claims.(IJWTExpire); ok && instance.GetExpiresAfter() != nil && now().UnixNano() > instance.GetExpiresAfter().Time().UnixNano() {
		return ErrTokenExpired
	}
	return nil
}
