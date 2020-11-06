package resource

import (
	"errors"
	"github.com/NavPool/navpool-hq-api/internal/service"
	"github.com/NavPool/navpool-hq-api/internal/service/account"
	"github.com/NavPool/navpool-hq-api/internal/service/twofactor"
	"github.com/gin-gonic/gin"
)

type TwoFactorResource struct {
	accounts *service.AccountService
	service  *service.TwoFactorService
}

func NewTwoFactorResource(accountService *service.AccountService, twoFactorService *service.TwoFactorService) *TwoFactorResource {
	return &TwoFactorResource{accountService, twoFactorService}
}

var (
	ErrTwoFactorUnableToSetup = errors.New("Unable to generate new 2FA secret")
)

func (r *TwoFactorResource) GetTwoFactorSecret(c *gin.Context) {
	userId, _ := c.Get("id")
	user, err := r.accounts.GetUserByClaim(userId.(account.User), "TwoFactor")
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	if user.TwoFactorActive() {
		err = service.ErrTwoFactorAlreadyActive
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	otp, err := r.service.GetSecret(user.Username, user)
	if err != nil {
		err = ErrTwoFactorUnableToSetup
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, otp)
}

func (r *TwoFactorResource) EnableTwoFactor(c *gin.Context) {
	userId, _ := c.Get("id")
	user, err := r.accounts.GetUserByClaim(userId.(account.User), "TwoFactor")
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	verification := twofactor.Verification{}
	if err := c.BindJSON(&verification); err != nil {
		return
	}

	err = r.service.Enable(verification, user)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, nil)
}

func (r *TwoFactorResource) DisableTwoFactor(c *gin.Context) {
	userId, _ := c.Get("id")
	user, err := r.accounts.GetUserByClaim(userId.(account.User), "TwoFactor")
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	verification := twofactor.Verification{}
	if err := c.BindJSON(&verification); err != nil {
		return
	}

	err = r.service.Disable(verification, user)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, nil)
}
