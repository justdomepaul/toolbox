package restful

import (
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/justdomepaul/toolbox/definition"
	"github.com/justdomepaul/toolbox/errorhandler"
	jwtTool "github.com/justdomepaul/toolbox/jwt"
	"strings"
)

func NewBasicGuardValidator(jwt jwtTool.IJWT) *BasicGuardValidator {
	return &BasicGuardValidator{
		jwt: jwt,
	}
}

type BasicGuardValidator struct {
	jwt jwtTool.IJWT
}

func (b *BasicGuardValidator) Verify(c *gin.Context, token string) error {
	commonClaims := jwtTool.NewCommon(jwtTool.NewClaimsBuilder().Build())
	if err := b.jwt.VerifyToken(token, commonClaims); err != nil {
		return errorhandler.NewErrAuthenticate(err)
	}

	for _, permission := range commonClaims.Permissions {
		if strings.HasPrefix(c.FullPath(), permission) || strings.HasPrefix(c.Request.RequestURI, permission) {
			c.Set(definition.AuthTokenKey, commonClaims)
			return nil
		}
	}
	return errorhandler.NewErrPermissionDeny(errors.New("no permission allowed to resource"))
}
