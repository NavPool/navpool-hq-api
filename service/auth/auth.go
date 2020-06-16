package auth

import (
	"github.com/NavPool/navpool-hq-api/config"
	"github.com/NavPool/navpool-hq-api/logger"
	"github.com/NavPool/navpool-hq-api/service/account"
	"github.com/NavPool/navpool-hq-api/service/twofactor"
	"github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
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

	user, err := account.GetUserByUsernamePassword(loginVals.Username, loginVals.Password, "TwoFactor")
	if err != nil {
		logger.LogError(err)
		return nil, jwt.ErrFailedAuthentication
	}

	if user.TwoFactor.Active {
		success, lastUsed, err := twofactor.Verify(user.TwoFactor.Secret, loginVals.TwoFactor, user.TwoFactor.LastUsed)
		if err != nil || success == false {
			if err != nil {
				logger.LogError(err)
			}
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
	_ = account.UpdateUser(user)

	return user, nil
}

func Authorizator(data interface{}, c *gin.Context) bool {
	if _, ok := data.(account.User); ok {
		return true
	}

	return false
}

func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": "Username, password or 2FA is incorrect",
	})
}
