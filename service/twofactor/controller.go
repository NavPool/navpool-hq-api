package twofactor

import (
	"errors"
	"github.com/NavPool/navpool-hq-api/service/account"
	"github.com/gin-gonic/gin"
)

type Controller struct{}

var (
	ErrTwoFactorAlreadyActive = errors.New("2FA already active on account")
	ErrTwoFactorNotActive     = errors.New("2FA is not active on account")
	ErrTwoFactorUnableToSetup = errors.New("Unable to generate new 2FA secret")
	ErrTwoFactorInvalidCode   = errors.New("Authentication code is invalid")
)

func (controller *Controller) GetTwoFactorSecret(c *gin.Context) {
	userId, _ := c.Get("id")
	user, err := account.GetUserByClaim(userId.(account.User), "TwoFactor")
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	if user.TwoFactorActive() {
		err = ErrTwoFactorAlreadyActive
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	otp, err := GetSecret(user.Username, user)
	if err != nil {
		err = ErrTwoFactorUnableToSetup
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, otp)
}

func (controller *Controller) EnableTwoFactor(c *gin.Context) {
	userId, _ := c.Get("id")
	user, err := account.GetUserByClaim(userId.(account.User), "TwoFactor")
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	verification := Verification{}
	if err := c.BindJSON(&verification); err != nil {
		return
	}

	err = Enable(verification, user)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, nil)
}

func (controller *Controller) DisableTwoFactor(c *gin.Context) {
	userId, _ := c.Get("id")
	user, err := account.GetUserByClaim(userId.(account.User), "TwoFactor")
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	verification := Verification{}
	if err := c.BindJSON(&verification); err != nil {
		return
	}

	err = Disable(verification, user)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, nil)
}
