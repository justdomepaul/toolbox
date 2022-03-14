package jwt

import (
	"github.com/square/go-jose/v3/jwt"
	"time"
)

// ClaimsOption interface
type ClaimsOption interface {
	Apply(*Common)
}

// WithSecret method
func WithSecret(secret string) ClaimsOption {
	return withSecret{secret: secret}
}

type withSecret struct {
	secret string
}

// Apply method
func (w withSecret) Apply(c *Common) {
	c.Secret = w.secret
}

// Common type
type Common struct {
	Secret string `json:"s,omitempty"`
	*jwt.Claims
}

func (c *Common) ExpiresAfter(d time.Duration) {
	c.Expiry = expiresAfter(d)
}

func (c *Common) GetExpiresAfter() *jwt.NumericDate {
	return c.Expiry
}

// NewCommon method
func NewCommon(claims *jwt.Claims, options ...ClaimsOption) *Common {
	common := &Common{
		Claims: claims,
	}

	for _, option := range options {
		option.Apply(common)
	}

	return common
}
