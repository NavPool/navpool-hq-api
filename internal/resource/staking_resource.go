package resource

import (
	"errors"
	"github.com/NavPool/navpool-hq-api/internal/service"
	"github.com/NavPool/navpool-hq-api/internal/service/account"
	"github.com/gin-gonic/gin"
)

type StakingResource struct {
	staking *service.StakingService
}

func NewStakingResource(stakingService *service.StakingService) *StakingResource {
	return &StakingResource{stakingService}
}

var (
	ErrorUnableToGetStakingReport = errors.New("Unable to retrieve staking report")
)

func (r *StakingResource) GetStakingRewardsForAccount(c *gin.Context) {
	userId, _ := c.Get("id")
	user := userId.(account.User)

	rewards, err := r.staking.GetStakingRewardsForUser(user)
	if err != nil {
		err = ErrorUnableToGetStakingReport
		_ = c.Error(ErrorUnableToGetStakingReport).SetType(gin.ErrorTypePublic)
	}

	c.JSON(200, rewards)
}
