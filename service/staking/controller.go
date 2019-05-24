package staking

import (
	"errors"
	"github.com/NavPool/navpool-hq-api/service/account"
	"github.com/gin-gonic/gin"
)

type Controller struct{}

var (
	ErrorUnableToGetStakingReport = errors.New("Unable to retrieve staking report")
)

func (controller *Controller) GetStakingRewardsForAccount(c *gin.Context) {
	userId, _ := c.Get("id")
	user := userId.(account.User)

	rewards, err := GetStakingRewardsForUser(user)
	if err != nil {
		err = ErrorUnableToGetStakingReport
		_ = c.Error(ErrorUnableToGetStakingReport).SetType(gin.ErrorTypePublic)
	}

	c.JSON(200, rewards)
}
