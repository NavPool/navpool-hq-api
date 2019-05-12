package communityFund

import (
	"github.com/NavPool/navpool-hq-api/service/account"
	"github.com/NavPool/navpool-hq-api/service/communityFund/model"
	"github.com/gin-gonic/gin"
)

type Controller struct{}

type VoteDto struct {
	Type   model.VoteType   `json:"type" binding:"required"`
	Hash   string           `json:"hash" binding:"required"`
	Choice model.VoteChoice `json:"vote" binding:"required"`
}

func (controller *Controller) GetProposalVotes(c *gin.Context) {
	userId, _ := c.Get("id")
	user := userId.(account.User)

	votes, err := GetProposalVotes(user)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, votes)
}

func (controller *Controller) UpdateProposalVotes(c *gin.Context) {
	userId, _ := c.Get("id")
	user := userId.(account.User)

	votesDto := make([]VoteDto, 0)
	if err := c.BindJSON(&votesDto); err != nil {
		return
	}

	err := UpdateProposalVotes(votesDto, user)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, votesDto)
}

func (controller *Controller) GetPaymentRequestVotes(c *gin.Context) {
	userId, _ := c.Get("id")
	user := userId.(account.User)

	votes, err := GetPaymentRequestVotes(user)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, votes)
}

func (controller *Controller) UpdatePaymentRequestVotes(c *gin.Context) {
	userId, _ := c.Get("id")
	user := userId.(account.User)

	votesDto := make([]VoteDto, 0)
	if err := c.BindJSON(&votesDto); err != nil {
		return
	}

	err := UpdatePaymentRequestVotes(votesDto, user)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, votesDto)
}
