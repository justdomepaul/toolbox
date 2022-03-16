package restful

import (
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/justdomepaul/toolbox/config"
	"github.com/justdomepaul/toolbox/errorhandler"
	"github.com/justdomepaul/toolbox/jwt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

type testGuarderValidator struct {
	mock.Mock
	GuarderValidator
}

func (t *testGuarderValidator) Verify(c *gin.Context, token string) error {
	args := t.Called(c, token)
	return args.Error(0)
}

type MiddlewareSuite struct {
	suite.Suite
	jwtOp config.JWT
	jwt   jwt.IJWT
}

func (suite *MiddlewareSuite) SetupSuite() {
	suite.jwtOp = config.JWT{
		EcdsaPrivateKey: `-----BEGIN EC PRIVATE KEY-----
MHcCAQEEIOChaSphj1MdLSxvU56h9vwmmpqdsQQF2alVwLKTj7dMoAoGCCqGSM49
AwEHoUQDQgAE7gMib5EUeW1An5VkkY4aU3xy+altlU3U0zn3FCO9Ffe/wwNUcUzp
XC9HWu76KhJnPpHczvZZv7Rro+kmqvN5tw==
-----END EC PRIVATE KEY-----
`}
	es256, err := jwt.NewES256JWT(suite.jwtOp.EcdsaPrivateKey)
	suite.NoError(err)
	suite.jwt = es256

}

func (suite *MiddlewareSuite) TestNewJWTGuard() {
	testGuarderValidator := &testGuarderValidator{}
	testGuarderValidator.On("Verify", mock.Anything, mock.Anything).Return(nil)
	result := NewJWTGuarder(suite.jwtOp, testGuarderValidator).JWTGuarder()
	suite.Equal("gin.HandlerFunc", reflect.TypeOf(result).String())
}

func (suite *MiddlewareSuite) TestJWTGuarderRun() {
	token, err := suite.jwt.GenerateToken(jwt.NewCommon(
		jwt.NewClaimsBuilder().ExpiresAfter(500*time.Second).Build(),
		jwt.WithPermissions("/ping"),
	))
	suite.NoError(err)

	testGuarderValidator := &testGuarderValidator{}
	testGuarderValidator.On("Verify", mock.Anything, token).Return(nil)

	r := gin.Default()
	r.GET("/ping", NewJWTGuarder(suite.jwtOp, testGuarderValidator).JWTGuarder())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
}

func (suite *MiddlewareSuite) TestJWTGuarderRunWhiteList() {
	token, err := suite.jwt.GenerateToken(jwt.NewCommon(
		jwt.NewClaimsBuilder().ExpiresAfter(500*time.Second).Build(),
		jwt.WithPermissions("/ping"),
	))
	suite.NoError(err)

	testGuarderValidator := &testGuarderValidator{}
	testGuarderValidator.On("Verify", mock.Anything, token).Return(nil)

	r := gin.Default()
	r.GET("/ping", NewJWTGuarder(suite.jwtOp, testGuarderValidator).JWTGuarder("/ping"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)
}

func (suite *MiddlewareSuite) TestJWTGuarderRunFormatError() {
	token, err := suite.jwt.GenerateToken(jwt.NewCommon(jwt.NewClaimsBuilder().Build()))
	suite.NoError(err)

	testGuarderValidator := &testGuarderValidator{}
	testGuarderValidator.On("Verify", mock.Anything, token).Return(nil)

	r := gin.Default()
	r.Use(gin.Logger(), errorhandler.GinPanicErrorHandler("Gin Mock", "Gin Mock test JWT guard"))
	r.GET("/ping", NewJWTGuarder(suite.jwtOp, testGuarderValidator).JWTGuarder())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Add("Authorization", "Basic "+token)
	r.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
}

func (suite *MiddlewareSuite) TestJWTGuarderRunExpired() {
	token, err := suite.jwt.GenerateToken(jwt.NewCommon(jwt.NewClaimsBuilder().ExpiresAfter(1 * time.Second).Build()))
	suite.NoError(err)

	time.Sleep(1 * time.Second)

	testGuarderValidator := &testGuarderValidator{}
	testGuarderValidator.On("Verify", mock.Anything, token).Return(errorhandler.NewErrPermissionDeny(errors.New("got error")))

	r := gin.Default()
	r.Use(gin.Logger(), errorhandler.GinPanicErrorHandler("Gin Mock", "Gin Mock test JWT guard"))
	r.GET("/ping", NewJWTGuarder(suite.jwtOp, testGuarderValidator).JWTGuarder())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	suite.Equal(http.StatusForbidden, w.Code)
}

func (suite *MiddlewareSuite) TestJWTGuarderRunPermissionNotAllow() {
	token, err := suite.jwt.GenerateToken(jwt.NewCommon(
		jwt.NewClaimsBuilder().ExpiresAfter(500 * time.Second).Build(),
	))
	suite.NoError(err)

	testGuarderValidator := &testGuarderValidator{}
	testGuarderValidator.On("Verify", mock.Anything, token).Return(errorhandler.NewErrPermissionDeny(errors.New("got error")))

	r := gin.Default()
	r.Use(gin.Logger(), errorhandler.GinPanicErrorHandler("Gin Mock", "Gin Mock test JWT guard"))
	r.GET("/ping", NewJWTGuarder(suite.jwtOp, testGuarderValidator).JWTGuarder())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Add("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	suite.Equal(http.StatusForbidden, w.Code)
}

func TestMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(MiddlewareSuite))
}
