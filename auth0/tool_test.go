package auth0

import (
	"crypto/rsa"
	"github.com/justdomepaul/toolbox/jwt"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"
	"math/big"
	"testing"
)

type Auth0Suite struct {
	suite.Suite
	auth0Domain string
	cert        string
	idToken     string
}

func (suite *Auth0Suite) SetupSuite() {
	suite.auth0Domain = "https://dev-te5750n303no4b6q.us.auth0.com"
	suite.cert = "MIIDHTCCAgWgAwIBAgIJCPhayegd+1zUMA0GCSqGSIb3DQEBCwUAMCwxKjAoBgNVBAMTIWRldi10ZTU3NTBuMzAzbm80YjZxLnVzLmF1dGgwLmNvbTAeFw0yMzAzMjkxNjU4MjlaFw0zNjEyMDUxNjU4MjlaMCwxKjAoBgNVBAMTIWRldi10ZTU3NTBuMzAzbm80YjZxLnVzLmF1dGgwLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKra4G9PH9D9VqE/d+zt/rP+oq19rAce/8Ml+zDDNYpZkWbBOX9O3zx7UIzoBAbg74MhN+xP0jbKocIqD86wAm7gtdCF7ZHU92ox0+Ef/Q1w69a7HvdaryGW9wQDT5nZtHVbVK0QxwN7iproEsxA22Fe1xVIoYrlbYcHxnz0WzMgAyApdx/M8nh4WAT8yhN+O2z3zLGvTXVVGI7m85wQny/H7p1QOMD5jeLPjql6og3m0MwpgI9HE9pfTvtk/eU+MVSYcseWocAw/olKHAiQb7yiZW0dDfo44t0nT7+JLGTz+/g6XKCUs0X4fJeFmr9x64MAvwLcZaIOqQaoe5gNlWsCAwEAAaNCMEAwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUtVk6RpxWqVDhtZ4y5vmgHofOTgMwDgYDVR0PAQH/BAQDAgKEMA0GCSqGSIb3DQEBCwUAA4IBAQBT0Cqdw7OjTG/tpzNYMwNdjpv+oD3Le4Fz11JgtXGHbinEiZ6T7pAkFoti6ee5b9CHUfOLBUGA0/tA9+kXxtR/dNGdt+ohmkLMY4TDZcQAEsBpz+JpnRDyAU2KMwywPMBGXC31t29ivTONC8SISEfVwyBdIExDI8IHl1iVcqNjee7N/ho4BHH49dEZqjSy+Bsxpay16Pe2kv9D80HeePLCXPN8UTk7BRh+0aXs+CmR4E3+xLze0BPH9YRpTqrDuGDfBmY7WlCUjd2T1FvLg0ZaNIMQw2dt6BWc8IkQvRLZT+gnN53z5JmjRGPuCWl84AOjvK3WqwgHxFOahfxNEDF9"
	suite.idToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjhqeXJ6R1RxWkpxMm4taHlCQUtKaSJ9.eyJ1c2VyX3Blcm1pc3Npb25zIjpbInJlYWQ6dG9kbyIsIndyaXRlOmNyZWF0ZVRvZG8iXSwibmlja25hbWUiOiJqdXN0ZG9tZXBhdWwiLCJuYW1lIjoiTWF4Zm9ja2VyIiwicGljdHVyZSI6Imh0dHBzOi8vYXZhdGFycy5naXRodWJ1c2VyY29udGVudC5jb20vdS8xMDE3NDkyND92PTQiLCJ1cGRhdGVkX2F0IjoiMjAyMy0wNC0wMlQxNzo0NToxMi4xMzVaIiwiZW1haWwiOiJqdXN0ZG9tZXBhdWxAZ21haWwuY29tIiwiaXNzIjoiaHR0cHM6Ly9kZXYtdGU1NzUwbjMwM25vNGI2cS51cy5hdXRoMC5jb20vIiwiYXVkIjoiYWcxTDJpMUVkbnBsSnRsR0Rkak9aMHhTc3BrWGRuTVAiLCJpYXQiOjE2ODA0NTc1MTMsImV4cCI6MTY4MDQ5MzUxMywic3ViIjoiZ2l0aHVifDEwMTc0OTI0IiwiYXV0aF90aW1lIjoxNjgwNDU3NTEyLCJzaWQiOiJpa3JzYlk0V2FERnBOWGhGdDFYMDZ2OXQwU0xXRHM1WSIsIm5vbmNlIjoiYmtsb1R6TldNazFaVjFaTmVrWTJlV3BYT0RaTWVWUk9UblJJV2xadlpXTlhZelJ0UVZCellWbGxRUT09In0.MvRtxYgy6ZkQB6Vi9pN717tAJddKf9l07Rt_NtLOWWrCRH2CQ8Gj3ddpIJVT0tg2dDWtzJYrQP-z8inHgz3Htx-0SgTuXDfcIWxM0g0cMQUM5FgVEWS73HFIDEz2HV8XBg1rsT1Duvyjp0FgHyli6ku_eX0SC48K2y0bD5Pp6GlIPlfceDGM-H5pl3YyW9JypKDwsboy308TRHl38xwzG88ppyihaQl6y6cdsOdZuF-i62F3Kgou_MmWiL0eN9WiOGI_xVADa6DURw8xaHHvoz4DtsnplL8NCAomdb4FUKF22Q7eNp_VtGNyNJsZKY4N6evZ5LSaPw8EvO0Z9sm2GA"
}

func (suite *Auth0Suite) TestGetAuth0JWKSInfo() {
	jwksInfo, err := GetAuth0JWKSInfo(suite.auth0Domain)
	suite.NoError(err)
	suite.NotEmpty(jwksInfo)
	suite.Greater(len(gjson.GetBytes(jwksInfo, "keys").Array()), 0)
}

func (suite *Auth0Suite) TestGetAuth0JWKSInfoFailDomain() {
	_, err := GetAuth0JWKSInfo("https://localhost")
	suite.Error(err)
}

func (suite *Auth0Suite) TestGetAuth0JWKSSet() {
	jwksInfo, err := GetAuth0JWKSInfo(suite.auth0Domain)
	suite.NoError(err)
	set, err := GetAuth0JWKSSet(jwksInfo)
	suite.NoError(err)
	suite.Greater(len(set.Keys), 0)
	suite.T().Logf("%+v\n", set)
}

func (suite *Auth0Suite) TestGetAuth0JWKSSetFailJWKSInfo() {
	_, err := GetAuth0JWKSSet([]byte(""))
	suite.Error(err)
}

func (suite *Auth0Suite) TestGetAuth0JWKSPublicKeyCert() {
	jwksInfo, err := GetAuth0JWKSInfo(suite.auth0Domain)
	suite.NoError(err)
	cert, err := GetAuth0JWKSPublicKeyCert(jwksInfo, suite.idToken)
	suite.NoError(err)
	suite.Equal(suite.cert, cert)
}

func (suite *Auth0Suite) TestParseAuth0RSAPublicKeyFromCert() {
	publicKey, err := ParseAuth0RSAPublicKeyFromCert(suite.cert)
	suite.NoError(err)
	suite.Equal(65537, publicKey.E)
	expectedN, ok := big.NewInt(0).SetString("21568443966916899517039010919720256730322266476936012457374985944419195548002372631345442991963484866777373696717000859586227789796668767799372635072435648714244570167067390450368018805362097163101674077021814767372854876987420483266401263705179361997301256452454045996541524487778864427302147042034211142499499628864978964327367166000186836188945367288494309079779965093512563359747870109765265801040762496100182975730637665438048683241556199153641032853060469277902907114880162362114799936294763416310455910343901843510294343389761050636005580142797157399607563984405262293618172282167538564302071225779359197271403", 10)
	suite.True(ok)
	suite.Equal(expectedN, publicKey.N)
}

func (suite *Auth0Suite) TestValidateAuth0RS256IDTokenByJWKSSet() {
	jwksInfo, err := GetAuth0JWKSInfo(suite.auth0Domain)
	suite.NoError(err)
	set, err := GetAuth0JWKSSet(jwksInfo)
	suite.NoError(err)
	suite.NoError(ValidateAuth0RS256IDTokenByJWKSSet(set, suite.idToken))
}

func (suite *Auth0Suite) TestValidateAuth0RS256IDTokenByJWKSSetFail() {
	jwksInfo, err := GetAuth0JWKSInfo(suite.auth0Domain)
	suite.NoError(err)
	set, err := GetAuth0JWKSSet(jwksInfo)
	suite.NoError(err)
	suite.Error(ValidateAuth0RS256IDTokenByJWKSSet(set, "test"))
}

func (suite *Auth0Suite) TestValidateAuth0RS256IDToken() {
	publicKey, err := ParseAuth0RSAPublicKeyFromCert(suite.cert)
	suite.NoError(err)
	suite.NoError(ValidateAuth0RS256IDToken(publicKey, suite.idToken))
}

func (suite *Auth0Suite) TestValidateAuth0RS256IDTokenFail() {
	suite.Error(ValidateAuth0RS256IDToken(&rsa.PublicKey{}, "test"))
}

func (suite *Auth0Suite) TestVerifyAuth0RS256IDTokenByJWKSSet() {
	jwksInfo, err := GetAuth0JWKSInfo(suite.auth0Domain)
	suite.NoError(err)
	set, err := GetAuth0JWKSSet(jwksInfo)
	suite.NoError(err)
	claims := NewCommon(jwt.NewClaimsBuilder().Build())
	suite.ErrorIs(VerifyAuth0RS256IDTokenByJWKSSet(set, suite.idToken, claims), jwt.ErrTokenExpired)

	suite.Equal("justdomepaul@gmail.com", claims.Email)
	suite.Equal("Maxfocker", claims.Name)
	suite.Equal("justdomepaul", claims.Nickname)
	suite.Equal("bkloTzNWMk1ZV1ZNekY2eWpXODZMeVROTnRIWlZvZWNXYzRtQVBzYVllQQ==", claims.Nonce)
	suite.Equal("ikrsbY4WaDFpNXhFt1X06v9t0SLWDs5Y", claims.SID)
}

func (suite *Auth0Suite) TestVerifyAuth0RS256IDTokenByJWKSSetFailIDToken() {
	jwksInfo, err := GetAuth0JWKSInfo(suite.auth0Domain)
	suite.NoError(err)
	set, err := GetAuth0JWKSSet(jwksInfo)
	suite.NoError(err)
	claims := NewCommon(jwt.NewClaimsBuilder().Build())
	suite.Error(VerifyAuth0RS256IDTokenByJWKSSet(set, "test", claims))
}

func (suite *Auth0Suite) TestVerifyAuth0RS256IDToken() {
	publicKey, err := ParseAuth0RSAPublicKeyFromCert(suite.cert)
	suite.NoError(err)
	claims := NewCommon(jwt.NewClaimsBuilder().Build())
	suite.ErrorIs(VerifyAuth0RS256IDToken(publicKey, suite.idToken, claims), jwt.ErrTokenExpired)

	suite.Equal("justdomepaul@gmail.com", claims.Email)
	suite.Equal("Maxfocker", claims.Name)
	suite.Equal("justdomepaul", claims.Nickname)
	suite.Equal("bkloTzNWMk1ZV1ZNekY2eWpXODZMeVROTnRIWlZvZWNXYzRtQVBzYVllQQ==", claims.Nonce)
	suite.Equal("ikrsbY4WaDFpNXhFt1X06v9t0SLWDs5Y", claims.SID)
}

func (suite *Auth0Suite) TestVerifyAuth0RS256IDTokenFailIDToken() {
	publicKey, err := ParseAuth0RSAPublicKeyFromCert(suite.cert)
	suite.NoError(err)
	claims := NewCommon(jwt.NewClaimsBuilder().Build())
	suite.Error(VerifyAuth0RS256IDToken(publicKey, "test", claims))
}

func TestAuth0Suite(t *testing.T) {
	suite.Run(t, new(Auth0Suite))
}

func BenchmarkGetAuth0JWKSPublicKeyCert(b *testing.B) {
	auth0Domain := "https://dev-te5750n303no4b6q.us.auth0.com"
	idToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjhqeXJ6R1RxWkpxMm4taHlCQUtKaSJ9.eyJ1c2VyX3Blcm1pc3Npb25zIjpbInJlYWQ6dG9kbyIsIndyaXRlOmNyZWF0ZVRvZG8iXSwibmlja25hbWUiOiJqdXN0ZG9tZXBhdWwiLCJuYW1lIjoiTWF4Zm9ja2VyIiwicGljdHVyZSI6Imh0dHBzOi8vYXZhdGFycy5naXRodWJ1c2VyY29udGVudC5jb20vdS8xMDE3NDkyND92PTQiLCJ1cGRhdGVkX2F0IjoiMjAyMy0wNC0wMlQxNzo0NToxMi4xMzVaIiwiZW1haWwiOiJqdXN0ZG9tZXBhdWxAZ21haWwuY29tIiwiaXNzIjoiaHR0cHM6Ly9kZXYtdGU1NzUwbjMwM25vNGI2cS51cy5hdXRoMC5jb20vIiwiYXVkIjoiYWcxTDJpMUVkbnBsSnRsR0Rkak9aMHhTc3BrWGRuTVAiLCJpYXQiOjE2ODA0NTc1MTMsImV4cCI6MTY4MDQ5MzUxMywic3ViIjoiZ2l0aHVifDEwMTc0OTI0IiwiYXV0aF90aW1lIjoxNjgwNDU3NTEyLCJzaWQiOiJpa3JzYlk0V2FERnBOWGhGdDFYMDZ2OXQwU0xXRHM1WSIsIm5vbmNlIjoiYmtsb1R6TldNazFaVjFaTmVrWTJlV3BYT0RaTWVWUk9UblJJV2xadlpXTlhZelJ0UVZCellWbGxRUT09In0.MvRtxYgy6ZkQB6Vi9pN717tAJddKf9l07Rt_NtLOWWrCRH2CQ8Gj3ddpIJVT0tg2dDWtzJYrQP-z8inHgz3Htx-0SgTuXDfcIWxM0g0cMQUM5FgVEWS73HFIDEz2HV8XBg1rsT1Duvyjp0FgHyli6ku_eX0SC48K2y0bD5Pp6GlIPlfceDGM-H5pl3YyW9JypKDwsboy308TRHl38xwzG88ppyihaQl6y6cdsOdZuF-i62F3Kgou_MmWiL0eN9WiOGI_xVADa6DURw8xaHHvoz4DtsnplL8NCAomdb4FUKF22Q7eNp_VtGNyNJsZKY4N6evZ5LSaPw8EvO0Z9sm2GA"
	jwksInfo, _ := GetAuth0JWKSInfo(auth0Domain)
	for i := 0; i < b.N; i++ {
		GetAuth0JWKSPublicKeyCert(jwksInfo, idToken)
	}
}

func BenchmarkParseAuth0RSAPublicKeyFromCert(b *testing.B) {
	cert := "MIIDHTCCAgWgAwIBAgIJCPhayegd+1zUMA0GCSqGSIb3DQEBCwUAMCwxKjAoBgNVBAMTIWRldi10ZTU3NTBuMzAzbm80YjZxLnVzLmF1dGgwLmNvbTAeFw0yMzAzMjkxNjU4MjlaFw0zNjEyMDUxNjU4MjlaMCwxKjAoBgNVBAMTIWRldi10ZTU3NTBuMzAzbm80YjZxLnVzLmF1dGgwLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKra4G9PH9D9VqE/d+zt/rP+oq19rAce/8Ml+zDDNYpZkWbBOX9O3zx7UIzoBAbg74MhN+xP0jbKocIqD86wAm7gtdCF7ZHU92ox0+Ef/Q1w69a7HvdaryGW9wQDT5nZtHVbVK0QxwN7iproEsxA22Fe1xVIoYrlbYcHxnz0WzMgAyApdx/M8nh4WAT8yhN+O2z3zLGvTXVVGI7m85wQny/H7p1QOMD5jeLPjql6og3m0MwpgI9HE9pfTvtk/eU+MVSYcseWocAw/olKHAiQb7yiZW0dDfo44t0nT7+JLGTz+/g6XKCUs0X4fJeFmr9x64MAvwLcZaIOqQaoe5gNlWsCAwEAAaNCMEAwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUtVk6RpxWqVDhtZ4y5vmgHofOTgMwDgYDVR0PAQH/BAQDAgKEMA0GCSqGSIb3DQEBCwUAA4IBAQBT0Cqdw7OjTG/tpzNYMwNdjpv+oD3Le4Fz11JgtXGHbinEiZ6T7pAkFoti6ee5b9CHUfOLBUGA0/tA9+kXxtR/dNGdt+ohmkLMY4TDZcQAEsBpz+JpnRDyAU2KMwywPMBGXC31t29ivTONC8SISEfVwyBdIExDI8IHl1iVcqNjee7N/ho4BHH49dEZqjSy+Bsxpay16Pe2kv9D80HeePLCXPN8UTk7BRh+0aXs+CmR4E3+xLze0BPH9YRpTqrDuGDfBmY7WlCUjd2T1FvLg0ZaNIMQw2dt6BWc8IkQvRLZT+gnN53z5JmjRGPuCWl84AOjvK3WqwgHxFOahfxNEDF9"
	for i := 0; i < b.N; i++ {
		ParseAuth0RSAPublicKeyFromCert(cert)
	}
}

func BenchmarkVerifyAuth0RS256IDToken(b *testing.B) {
	cert := "MIIDHTCCAgWgAwIBAgIJCPhayegd+1zUMA0GCSqGSIb3DQEBCwUAMCwxKjAoBgNVBAMTIWRldi10ZTU3NTBuMzAzbm80YjZxLnVzLmF1dGgwLmNvbTAeFw0yMzAzMjkxNjU4MjlaFw0zNjEyMDUxNjU4MjlaMCwxKjAoBgNVBAMTIWRldi10ZTU3NTBuMzAzbm80YjZxLnVzLmF1dGgwLmNvbTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAKra4G9PH9D9VqE/d+zt/rP+oq19rAce/8Ml+zDDNYpZkWbBOX9O3zx7UIzoBAbg74MhN+xP0jbKocIqD86wAm7gtdCF7ZHU92ox0+Ef/Q1w69a7HvdaryGW9wQDT5nZtHVbVK0QxwN7iproEsxA22Fe1xVIoYrlbYcHxnz0WzMgAyApdx/M8nh4WAT8yhN+O2z3zLGvTXVVGI7m85wQny/H7p1QOMD5jeLPjql6og3m0MwpgI9HE9pfTvtk/eU+MVSYcseWocAw/olKHAiQb7yiZW0dDfo44t0nT7+JLGTz+/g6XKCUs0X4fJeFmr9x64MAvwLcZaIOqQaoe5gNlWsCAwEAAaNCMEAwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUtVk6RpxWqVDhtZ4y5vmgHofOTgMwDgYDVR0PAQH/BAQDAgKEMA0GCSqGSIb3DQEBCwUAA4IBAQBT0Cqdw7OjTG/tpzNYMwNdjpv+oD3Le4Fz11JgtXGHbinEiZ6T7pAkFoti6ee5b9CHUfOLBUGA0/tA9+kXxtR/dNGdt+ohmkLMY4TDZcQAEsBpz+JpnRDyAU2KMwywPMBGXC31t29ivTONC8SISEfVwyBdIExDI8IHl1iVcqNjee7N/ho4BHH49dEZqjSy+Bsxpay16Pe2kv9D80HeePLCXPN8UTk7BRh+0aXs+CmR4E3+xLze0BPH9YRpTqrDuGDfBmY7WlCUjd2T1FvLg0ZaNIMQw2dt6BWc8IkQvRLZT+gnN53z5JmjRGPuCWl84AOjvK3WqwgHxFOahfxNEDF9"
	idToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjhqeXJ6R1RxWkpxMm4taHlCQUtKaSJ9.eyJ1c2VyX3Blcm1pc3Npb25zIjpbInJlYWQ6dG9kbyIsIndyaXRlOmNyZWF0ZVRvZG8iXSwibmlja25hbWUiOiJqdXN0ZG9tZXBhdWwiLCJuYW1lIjoiTWF4Zm9ja2VyIiwicGljdHVyZSI6Imh0dHBzOi8vYXZhdGFycy5naXRodWJ1c2VyY29udGVudC5jb20vdS8xMDE3NDkyND92PTQiLCJ1cGRhdGVkX2F0IjoiMjAyMy0wNC0wMlQxNzo0NToxMi4xMzVaIiwiZW1haWwiOiJqdXN0ZG9tZXBhdWxAZ21haWwuY29tIiwiaXNzIjoiaHR0cHM6Ly9kZXYtdGU1NzUwbjMwM25vNGI2cS51cy5hdXRoMC5jb20vIiwiYXVkIjoiYWcxTDJpMUVkbnBsSnRsR0Rkak9aMHhTc3BrWGRuTVAiLCJpYXQiOjE2ODA0NTc1MTMsImV4cCI6MTY4MDQ5MzUxMywic3ViIjoiZ2l0aHVifDEwMTc0OTI0IiwiYXV0aF90aW1lIjoxNjgwNDU3NTEyLCJzaWQiOiJpa3JzYlk0V2FERnBOWGhGdDFYMDZ2OXQwU0xXRHM1WSIsIm5vbmNlIjoiYmtsb1R6TldNazFaVjFaTmVrWTJlV3BYT0RaTWVWUk9UblJJV2xadlpXTlhZelJ0UVZCellWbGxRUT09In0.MvRtxYgy6ZkQB6Vi9pN717tAJddKf9l07Rt_NtLOWWrCRH2CQ8Gj3ddpIJVT0tg2dDWtzJYrQP-z8inHgz3Htx-0SgTuXDfcIWxM0g0cMQUM5FgVEWS73HFIDEz2HV8XBg1rsT1Duvyjp0FgHyli6ku_eX0SC48K2y0bD5Pp6GlIPlfceDGM-H5pl3YyW9JypKDwsboy308TRHl38xwzG88ppyihaQl6y6cdsOdZuF-i62F3Kgou_MmWiL0eN9WiOGI_xVADa6DURw8xaHHvoz4DtsnplL8NCAomdb4FUKF22Q7eNp_VtGNyNJsZKY4N6evZ5LSaPw8EvO0Z9sm2GA"
	publicKey, _ := ParseAuth0RSAPublicKeyFromCert(cert)
	for i := 0; i < b.N; i++ {
		ValidateAuth0RS256IDToken(publicKey, idToken)
	}
}
