package jwt

import (
	"github.com/go-jose/go-jose/v3/jwt"
	"time"
)

var (
	// Now variable
	Now = time.Now
)

type IJWTClaims interface {
	Validate(e jwt.Expected) error
	ValidateWithLeeway(e jwt.Expected, leeway time.Duration) error
}

type IJWTExpire interface {
	ExpiresAfter(d time.Duration)
	GetExpiresAfter() *jwt.NumericDate
}

type IClaims interface {
	ExpiresAfter(d time.Duration) *ClaimsBuilder
	GetExpiresAfter() *jwt.NumericDate
	WithAudience(audience jwt.Audience) *ClaimsBuilder
	GetAudience() jwt.Audience
	WithID(id string) *ClaimsBuilder
	GetID() string
	WithIssuedAt() *ClaimsBuilder
	GetIssuedAt() *jwt.NumericDate
	WithIssuer(issuer string) *ClaimsBuilder
	GetIssuer() string
	NotUseBefore(d time.Duration) *ClaimsBuilder
	GetNotBefore() *jwt.NumericDate
	WithSubject(subject string) *ClaimsBuilder
	GetSubject() string
	Build() *jwt.Claims
}

func expiresAfter(d time.Duration) *jwt.NumericDate {
	return jwt.NewNumericDate(Now().Add(d))
}

// ClaimsBuilder type
type ClaimsBuilder jwt.Claims

// WithAudience method
func (c *ClaimsBuilder) WithAudience(audience jwt.Audience) *ClaimsBuilder {
	c.Audience = audience
	return c
}

// GetAudience method
func (c *ClaimsBuilder) GetAudience() jwt.Audience {
	return c.Audience
}

// ExpiresAfter method
func (c *ClaimsBuilder) ExpiresAfter(d time.Duration) *ClaimsBuilder {
	c.Expiry = expiresAfter(d)
	return c
}

// GetExpiresAfter method
func (c *ClaimsBuilder) GetExpiresAfter() *jwt.NumericDate {
	return c.Expiry
}

// WithID method
func (c *ClaimsBuilder) WithID(id string) *ClaimsBuilder {
	c.ID = id
	return c
}

// GetID method
func (c *ClaimsBuilder) GetID() string {
	return c.ID
}

// WithIssuedAt method
func (c *ClaimsBuilder) WithIssuedAt() *ClaimsBuilder {
	c.IssuedAt = jwt.NewNumericDate(Now())
	return c
}

// GetIssuedAt method
func (c *ClaimsBuilder) GetIssuedAt() *jwt.NumericDate {
	return c.IssuedAt
}

// WithIssuer method
func (c *ClaimsBuilder) WithIssuer(issuer string) *ClaimsBuilder {
	c.Issuer = issuer
	return c
}

// GetIssuer method
func (c *ClaimsBuilder) GetIssuer() string {
	return c.Issuer
}

// NotUseBefore method
func (c *ClaimsBuilder) NotUseBefore(d time.Duration) *ClaimsBuilder {
	c.NotBefore = jwt.NewNumericDate(Now().Add(-d))
	return c
}

// GetNotBefore method
func (c *ClaimsBuilder) GetNotBefore() *jwt.NumericDate {
	return c.NotBefore
}

// WithSubject method
func (c *ClaimsBuilder) WithSubject(subject string) *ClaimsBuilder {
	c.Subject = subject
	return c
}

// GetSubject method
func (c *ClaimsBuilder) GetSubject() string {
	return c.Subject
}

// Build method
func (c *ClaimsBuilder) Build() *jwt.Claims {
	// we could check everything here
	return (*jwt.Claims)(c)
}

// NewClaimsBuilder method
func NewClaimsBuilder() *ClaimsBuilder {
	return &ClaimsBuilder{}
}
