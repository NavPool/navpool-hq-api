package auth

import (
	"github.com/NavPool/navpool-hq-api/config"
	"github.com/NavPool/navpool-hq-api/service/account"
	"github.com/NavPool/navpool-hq-api/service/twofactor"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"log"
	"time"
)

var identityKey = config.Get().JWT.IdentityKey

func Middleware() (*jwt.GinJWTMiddleware, error) {
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           config.Get().JWT.Realm,
		Key:             []byte(config.Get().JWT.Secret),
		Timeout:         time.Hour,
		MaxRefresh:      time.Hour,
		IdentityKey:     identityKey,
		PayloadFunc:     Payload,
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

	log.Println("Couldn't map the payload to a user")

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
		log.Printf("Username or Password incorrect for %s", loginVals.Username)
		return nil, jwt.ErrFailedAuthentication
	}

	if user.TwoFactor.Active {
		log.Printf("2fa active for %s", loginVals.Username)
		success, lastUsed, err := twofactor.Verify(user.TwoFactor.Secret, loginVals.TwoFactor, user.TwoFactor.LastUsed)
		if err != nil || success == false {
			return nil, jwt.ErrFailedAuthentication
		}

		user.TwoFactor.LastUsed = &lastUsed
		lastLoginAt := time.Now().UTC()
		user.LastLoginAt = &lastLoginAt
		account.UpdateUser(user)
	}

	return user, nil
}

func Authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(account.User); ok {
		log.Printf("Authorized access for: %s", v.ID)
		return true
	}

	return false
}

func Unauthorized(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}

//func verify(user account.User, code string) (success bool, err error) {
//	var lastUsedList []int
//
//	otpConfig := &dgoogauth.OTPConfig{
//		Secret:      *user.TwoFactor.Secret,
//		WindowSize:  3,
//		DisallowReuse: append(lastUsedList, *user.TwoFactor.LastUsed),
//		HotpCounter: 0,
//		UTC: true,
//	}
//
//	return otpConfig.Authenticate(code)
//}
