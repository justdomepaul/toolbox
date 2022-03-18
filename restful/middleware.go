package restful

import (
	"github.com/gin-gonic/gin"
	"github.com/justdomepaul/toolbox/config"
	"github.com/justdomepaul/toolbox/definition"
	"github.com/justdomepaul/toolbox/errorhandler"
	"strings"
)

type GuarderValidator interface {
	Verify(c *gin.Context, token string) error
}

func NewJWTGuarder(option config.JWT, validator GuarderValidator) *JWTGuarder {
	return &JWTGuarder{
		option:    option,
		validator: validator,
	}
}

type JWTGuarder struct {
	option    config.JWT
	validator GuarderValidator
}

// JWTGuarder method
func (j *JWTGuarder) JWTGuarder(whitelist ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// skip if matching whitelist
		for _, term := range whitelist {
			if strings.HasPrefix(c.FullPath(), term) || strings.HasPrefix(c.Request.RequestURI, term) {
				c.Next()
				return
			}
		}
		token := c.DefaultQuery(definition.QueryAuthKey, "")
		if token == "" {
			authorization := c.GetHeader(definition.AuthorizationKey)
			if len(authorization) == 0 {
				panic(errorhandler.NewErrInvalidArgument(errorhandler.ErrAuthorizationRequired))
			}
			if strings.HasPrefix(authorization, definition.AuthorizationType) {
				token = strings.TrimPrefix(authorization, definition.AuthorizationType)
			} else {
				panic(errorhandler.NewErrInvalidArgument(errorhandler.ErrAuthorizationTypeBearer))
			}
		}

		if err := j.validator.Verify(c, token); err != nil {
			panic(err)
		}
	}
}
