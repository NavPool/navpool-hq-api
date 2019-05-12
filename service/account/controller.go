package account

import (
	"errors"
	"github.com/gin-gonic/gin"
)

type Controller struct{}

var (
	ErrUnableToValidateUser = errors.New("Unable to validate user")
)

func (controller *Controller) GetAccount(c *gin.Context) {
	userId, _ := c.Get("id")
	user, err := GetUserByClaim(userId.(User), "TwoFactor")
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, user)
}
