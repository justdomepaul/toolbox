package jwt

import (
	"github.com/prashantv/gostub"
	"github.com/square/go-jose/v3/jwt"
	"github.com/stretchr/testify/suite"
	"reflect"
	"testing"
	"time"
)

type ClaimsSuite struct {
	suite.Suite
}

func (suite *ClaimsSuite) TestNewClaimsBuilder() {
	suite.Equal("*jwt.ClaimsBuilder", reflect.TypeOf(NewClaimsBuilder()).String())
}

func (suite *ClaimsSuite) TestClaimsBuilderWithAudienceMethod() {
	suite.Equal(jwt.Audience{"test"}, NewClaimsBuilder().WithAudience([]string{"test"}).Audience)
}

func (suite *ClaimsSuite) TestClaimsBuilderGetAudienceMethod() {
	suite.Equal(jwt.Audience{"test"}, NewClaimsBuilder().WithAudience([]string{"test"}).GetAudience())
}

func (suite *ClaimsSuite) TestClaimsBuilderExpiresAfterMethod() {
	defer gostub.Stub(&Now, func() time.Time {
		return time.Date(2019, 02, 22, 11, 55, 20, 0, time.UTC)
	}).Reset()
	suite.Equal("*jwt.NumericDate", reflect.TypeOf(NewClaimsBuilder().ExpiresAfter(5*time.Second).Expiry).String())
}

func (suite *ClaimsSuite) TestClaimsBuilderGetExpiresAfterMethod() {
	defer gostub.Stub(&Now, func() time.Time {
		return time.Date(2019, 02, 22, 11, 55, 20, 0, time.UTC)
	}).Reset()
	suite.Equal("*jwt.NumericDate", reflect.TypeOf(NewClaimsBuilder().ExpiresAfter(5*time.Second).GetExpiresAfter()).String())
}

func (suite *ClaimsSuite) TestClaimsBuilderWithIDMethod() {
	suite.Equal("test", NewClaimsBuilder().WithID("test").ID)
}

func (suite *ClaimsSuite) TestClaimsBuilderGetIDMethod() {
	suite.Equal("test", NewClaimsBuilder().WithID("test").GetID())
}

func (suite *ClaimsSuite) TestClaimsBuilderWithIssuedAtMethod() {
	defer gostub.Stub(&Now, func() time.Time {
		return time.Date(2019, 02, 22, 11, 55, 20, 0, time.UTC)
	}).Reset()
	suite.Equal("*jwt.NumericDate", reflect.TypeOf(NewClaimsBuilder().WithIssuedAt().IssuedAt).String())
}

func (suite *ClaimsSuite) TestClaimsBuilderGetIssuedAtMethod() {
	defer gostub.Stub(&Now, func() time.Time {
		return time.Date(2019, 02, 22, 11, 55, 20, 0, time.UTC)
	}).Reset()
	suite.Equal("*jwt.NumericDate", reflect.TypeOf(NewClaimsBuilder().WithIssuedAt().GetIssuedAt()).String())
}

func (suite *ClaimsSuite) TestClaimsBuilderWithIssuerMethod() {
	suite.Equal("test", NewClaimsBuilder().WithIssuer("test").Issuer)
}

func (suite *ClaimsSuite) TestClaimsBuilderGetIssuerMethod() {
	suite.Equal("test", NewClaimsBuilder().WithIssuer("test").GetIssuer())
}

func (suite *ClaimsSuite) TestClaimsBuilderNotUseBeforeMethod() {
	defer gostub.Stub(&Now, func() time.Time {
		return time.Date(2019, 02, 22, 11, 55, 20, 0, time.UTC)
	}).Reset()
	suite.Equal("*jwt.NumericDate", reflect.TypeOf(NewClaimsBuilder().NotUseBefore(5*time.Second).NotBefore).String())
}

func (suite *ClaimsSuite) TestClaimsBuilderGetNotBeforeMethod() {
	defer gostub.Stub(&Now, func() time.Time {
		return time.Date(2019, 02, 22, 11, 55, 20, 0, time.UTC)
	}).Reset()
	suite.Equal("*jwt.NumericDate", reflect.TypeOf(NewClaimsBuilder().NotUseBefore(5*time.Second).GetNotBefore()).String())
}

func (suite *ClaimsSuite) TestClaimsBuilderWithSubjectMethod() {
	suite.Equal("test", NewClaimsBuilder().WithSubject("test").Subject)
}

func (suite *ClaimsSuite) TestClaimsBuilderGetSubjectMethod() {
	suite.Equal("test", NewClaimsBuilder().WithSubject("test").GetSubject())
}

func (suite *ClaimsSuite) TestClaimsBuilderBuildMethod() {
	claims := NewClaimsBuilder().Build()
	suite.Equal("*jwt.Claims", reflect.TypeOf(claims).String())
}

func TestClaimsSuite(t *testing.T) {
	suite.Run(t, new(ClaimsSuite))
}
