package jwt

import (
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/justdomepaul/toolbox/generic"
	"time"
)

// ClaimsOption interface
type ClaimsOption interface {
	Apply(*Common)
}

// WithRootID method
func WithRootID[T generic.ByteSeq](id T) ClaimsOption {
	return withRootID{id: []byte(id)}
}

type withRootID struct {
	id []byte
}

// Apply method
func (w withRootID) Apply(c *Common) {
	c.ClientID = w.id
}

// WithClientID method
func WithClientID[T generic.ByteSeq](id T) ClaimsOption {
	return withClientID{id: []byte(id)}
}

type withClientID struct {
	id []byte
}

// Apply method
func (w withClientID) Apply(c *Common) {
	c.ClientID = w.id
}

// WithSecret method
func WithSecret[T generic.ByteSeq](secret T) ClaimsOption {
	return withSecret{secret: []byte(secret)}
}

type withSecret struct {
	secret []byte
}

// Apply method
func (w withSecret) Apply(c *Common) {
	c.Secret = w.secret
}

// WithPermissions method
func WithPermissions(permissions ...string) ClaimsOption {
	return withPermissions{permissions: permissions}
}

type withPermissions struct {
	permissions []string
}

// Apply method
func (w withPermissions) Apply(t *Common) {
	t.Permissions = w.permissions
}

// WithScopes method
func WithScopes(scopes ...string) ClaimsOption {
	return withScopes{scopes: scopes}
}

type withScopes struct {
	scopes []string
}

// Apply method
func (w withScopes) Apply(t *Common) {
	t.Scopes = w.scopes
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

// Common type
type Common struct {
	Secret      []byte   `json:"s,omitempty"`
	RootID      []byte   `json:"root_id,omitempty"`
	ClientID    []byte   `json:"client_id,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
	Scopes      []string `json:"scopes,omitempty"`
	*jwt.Claims
}

func (c *Common) ExpiresAfter(d time.Duration) {
	c.Expiry = expiresAfter(d)
}

func (c *Common) GetExpiresAfter() *jwt.NumericDate {
	return c.Expiry
}
