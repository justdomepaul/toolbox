package auth0

import (
	"github.com/go-jose/go-jose/v3/jwt"
	"time"
)

var (
	// Now variable
	now = time.Now
)

// NewCommon method
func NewCommon(claims *jwt.Claims) *Common {
	common := &Common{
		Claims: claims,
	}
	return common
}

type Common struct {
	AuthTime        *jwt.NumericDate `json:"auth_time,omitempty"`
	Email           string           `json:"email,omitempty"`
	Name            string           `json:"name,omitempty"`
	Nickname        string           `json:"nickname,omitempty"`
	Nonce           string           `json:"nonce,omitempty"`
	Picture         string           `json:"picture,omitempty"`
	SID             string           `json:"sid,omitempty"`
	UpdatedAt       string           `json:"updated_at,omitempty"`
	UserPermissions []string         `json:"user_permissions,omitempty"`
	*jwt.Claims
}

func (c *Common) ExpiresAfter(d time.Duration) {
	c.Expiry = jwt.NewNumericDate(now().Add(d))
}

func (c *Common) GetExpiresAfter() *jwt.NumericDate {
	return c.Expiry
}
