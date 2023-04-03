package auth0

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"github.com/go-jose/go-jose/v3/jwt"
	jwtPkg "github.com/golang-jwt/jwt"
	jwtTool "github.com/justdomepaul/toolbox/jwt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func GetAuth0JWKSInfo(auth0Domain string) ([]byte, error) {
	resp, err := http.Get(auth0Domain + "/.well-known/jwks.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func GetAuth0JWKSPublicKeyCert(jwksInfo []byte, idToken string) (cert string, err error) {
	jwtHeader, err := base64.StdEncoding.DecodeString(strings.SplitN(idToken, ".", 2)[0])
	if err != nil {
		return "", err
	}
	for _, key := range gjson.GetBytes(jwksInfo, "keys").Array() {
		if key.Get("kid").String() == gjson.GetBytes(jwtHeader, "kid").String() {
			cert = key.Get("x5c").Array()[0].String()
		}
	}
	return
}

func ParseAuth0RSAPublicKeyFromCert(cert string) (*rsa.PublicKey, error) {
	return jwtPkg.ParseRSAPublicKeyFromPEM(
		[]byte(
			fmt.Sprintf("-----BEGIN CERTIFICATE-----\n%s\n-----END CERTIFICATE-----\n",
				strings.Join(regexp.MustCompile(`.{1,64}`).FindAllString(cert, -1), "\n"),
			),
		),
	)
}

func ValidateAuth0RS256IDToken(rsaPublicKey *rsa.PublicKey, idToken string) error {
	token, err := jwt.ParseSigned(idToken)
	if err != nil {
		return err
	}
	return token.Claims(rsaPublicKey)
}

func VerifyAuth0RS256IDToken(rsaPublicKey *rsa.PublicKey, idToken string, claims jwtTool.IJWTClaims) error {
	token, err := jwt.ParseSigned(idToken)
	if err != nil {
		return err
	}

	if err := token.Claims(rsaPublicKey, claims); err != nil {
		return err
	}
	return checkExpire(claims)
}

func checkExpire(claims jwtTool.IJWTClaims) error {
	if instance, ok := claims.(jwtTool.IJWTExpire); ok && instance.GetExpiresAfter() != nil && now().UnixNano() > instance.GetExpiresAfter().Time().UnixNano() {
		return jwtTool.ErrTokenExpired
	}
	return nil
}
