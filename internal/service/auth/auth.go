package auth

import (
	"github.com/NavPool/navpool-hq-api/internal/config"
	"github.com/NavPool/navpool-hq-api/internal/di"
	"github.com/NavPool/navpool-hq-api/internal/service/account"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"time"
)

var identityKey = config.Get().JWT.IdentityKey

func Middleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           config.Get().JWT.Realm,
		Key:             []byte(config.Get().JWT.Secret),
		Timeout:         time.Hour * 6,
		MaxRefresh:      time.Hour * 24,
		PayloadFunc:     Payload,
		IdentityKey:     identityKey,
		IdentityHandler: IdentityHandler,
		Authenticator:   Authenticator,
		Authorizator:    Authorizator,
		Unauthorized:    Unauthorized,
		TokenLookup:     "header: Authorization",
		TokenHeadName:   "Bearer",
		TimeFunc:        time.Now,
	})
}

func Payload(data interface{}) jwt.MapClaims {
	if v, ok := data.(account.User); ok {
		return jwt.MapClaims{
			identityKey: v.ID,
		}
	}

	return jwt.MapClaims{}
}

func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return account.User{
		ID: uuid.FromStringOrNil(claims["id"].(string)),
	}
}

func Authenticator(c *gin.Context) (interface{}, error) {
	var loginVals account.Login
	if err := c.ShouldBind(&loginVals); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}

	user, err := di.Get().GetAccountService().GetUserByUsernamePassword(loginVals.Username, loginVals.Password, "TwoFactor")
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	if user.TwoFactor.Active {
		success, lastUsed, err := di.Get().GetTwofactorService().Verify(user.TwoFactor.Secret, loginVals.TwoFactor, user.TwoFactor.LastUsed)
		if err != nil || success == false {
			return nil, jwt.ErrFailedAuthentication
		}

		user.TwoFactor.LastUsed = &lastUsed
	} else {
		if loginVals.TwoFactor != "" {
			return nil, jwt.ErrFailedAuthentication
		}
	}

	lastLoginAt := time.Now().UTC()
	user.LastLoginAt = &lastLoginAt
	_ = di.Get().GetAccountService().UpdateUser(user)

	return user, nil
}

func Authorizator(data interface{}, c *gin.Context) bool {
	if _, ok := data.(account.User); ok {
		return true
	}

	return false
}

func Unauthorized(c *gin.Context, code int, message string) {
	logrus.Errorf("Unauthorized: %s", message)
	c.JSON(code, gin.H{
		"code":    code,
		"message": "Username, password or 2FA is incorrect",
	})
}
