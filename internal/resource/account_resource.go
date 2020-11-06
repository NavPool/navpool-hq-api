package resource

import (
	"github.com/NavPool/navpool-hq-api/internal/service"
	"github.com/NavPool/navpool-hq-api/internal/service/account"
	"github.com/gin-gonic/gin"
)

type AccountResource struct {
	accounts *service.AccountService
}

func NewAccountResource(accountService *service.AccountService) *AccountResource {
	return &AccountResource{accountService}
}

func (r *AccountResource) GetAccount(c *gin.Context) {
	userId, _ := c.Get("id")
	user, err := r.accounts.GetUserByClaim(userId.(account.User), "TwoFactor")
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, user)
}
