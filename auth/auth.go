package auth

import (
	"github.com/NavPool/navpool-hq-api/config"
	"github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
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
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value.
		// This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
}

func Payload(data interface{}) jwt.MapClaims {
	if v, ok := data.(*User); ok {
		return jwt.MapClaims{
			identityKey: v.Username,
		}
	}

	return jwt.MapClaims{}
}

func IdentityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	return &User{
		Username: claims["id"].(string),
	}
}

func Authenticator(c *gin.Context) (user interface{}, err error) {
	var loginVals Login
	if err := c.ShouldBind(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}

	user, err = GetUserByUsernamePassword(loginVals.Username, loginVals.Password)
	if err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	return user, err
}

func Authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(*User); ok && v.Active == true {
		log.Print(data)
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
