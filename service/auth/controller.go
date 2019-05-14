package auth

import (
	"errors"
	"github.com/NavPool/navpool-hq-api/service/account"
	"github.com/gin-gonic/gin"
)

type Controller struct{}

var (
	ErrMissingRegisterValues       = errors.New("Username, password or password confirmation not provided")
	ErrPasswordsDontMatch          = errors.New("Passwords do not match")
	ErrPasswordtooShort            = errors.New("Password must be at least 8 characters")
	ErrUsernameAlreadyInUse        = errors.New("The username is already in use")
	ErrUserRegistrationUnavailable = errors.New("User registering is unavailable")
)

func (controller *Controller) Register(c *gin.Context) {
	var registerVals account.Register
	if err := c.ShouldBind(&registerVals); err != nil {
		_ = c.Error(ErrMissingRegisterValues).SetType(gin.ErrorTypePublic)
		return
	}

	if registerVals.Password != registerVals.PasswordConfirm {
		_ = c.Error(ErrPasswordsDontMatch).SetType(gin.ErrorTypePublic)
		return
	}

	if len(registerVals.Password) < 8 {
		_ = c.Error(ErrPasswordtooShort).SetType(gin.ErrorTypePublic)
		return
	}

	exists, err := account.UsernameExists(registerVals.Username)
	if err != nil || exists == true {
		_ = c.Error(ErrUsernameAlreadyInUse).SetType(gin.ErrorTypePublic)
		return
	}

	user, err := account.CreateUser(registerVals.Username, registerVals.Password)
	if err != nil {
		_ = c.Error(ErrUserRegistrationUnavailable).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, user)
}
