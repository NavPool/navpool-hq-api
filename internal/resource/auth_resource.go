package resource

import (
	"errors"
	"github.com/NavPool/navpool-hq-api/internal/service"
	"github.com/NavPool/navpool-hq-api/internal/service/account"
	"github.com/getsentry/raven-go"
	"github.com/gin-gonic/gin"
)

type AuthResource struct {
	accounts *service.AccountService
}

var (
	ErrMissingRegisterValues       = errors.New("Username, password or password confirmation not provided")
	ErrPasswordsDontMatch          = errors.New("Passwords do not match")
	ErrUsernameTooShort            = errors.New("Username must be at least 6 characters")
	ErrPasswordTooShort            = errors.New("Password must be at least 6 characters")
	ErrUsernameAlreadyInUse        = errors.New("The username is already in use")
	ErrUserRegistrationUnavailable = errors.New("User registering is unavailable")
)

func (r *AuthResource) Register(c *gin.Context) {
	var registerVals account.Register
	if err := c.ShouldBind(&registerVals); err != nil {
		raven.CaptureErrorAndWait(err, nil)
		_ = c.Error(ErrMissingRegisterValues).SetType(gin.ErrorTypePublic)
		return
	}

	if len(registerVals.Username) < 6 {
		_ = c.Error(ErrUsernameTooShort).SetType(gin.ErrorTypePublic)
		return
	}

	if registerVals.Password != registerVals.PasswordConfirm {
		_ = c.Error(ErrPasswordsDontMatch).SetType(gin.ErrorTypePublic)
		return
	}

	if len(registerVals.Password) < 6 {
		_ = c.Error(ErrPasswordTooShort).SetType(gin.ErrorTypePublic)
		return
	}

	exists, err := r.accounts.UsernameExists(registerVals.Username)
	if err != nil || exists == true {
		if err != nil {
			raven.CaptureErrorAndWait(err, nil)
		}
		_ = c.Error(ErrUsernameAlreadyInUse).SetType(gin.ErrorTypePublic)
		return
	}

	user, err := r.accounts.CreateUser(registerVals.Username, registerVals.Password)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		_ = c.Error(ErrUserRegistrationUnavailable).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, user)
}
