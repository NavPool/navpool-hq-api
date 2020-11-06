package resource

import (
	"github.com/NavPool/navpool-hq-api/internal/resource/dto"
	"github.com/NavPool/navpool-hq-api/internal/service"
	"github.com/NavPool/navpool-hq-api/internal/service/account"
	"github.com/gin-gonic/gin"
)

type DaoResource struct {
	dao *service.DaoService
}

func NewDaoResource(daoService *service.DaoService) *DaoResource {
	return &DaoResource{daoService}
}

func (r *DaoResource) GetProposalVotes(c *gin.Context) {
	userId, _ := c.Get("id")
	user := userId.(account.User)

	votes, err := r.dao.GetProposalVotes(user)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, votes)
}

func (r *DaoResource) UpdateProposalVotes(c *gin.Context) {
	userId, _ := c.Get("id")
	user := userId.(account.User)

	votesDto := make([]*dto.Vote, 0)
	if err := c.BindJSON(&votesDto); err != nil {
		return
	}

	err := r.dao.UpdateProposalVotes(votesDto, user)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, votesDto)
}

func (r *DaoResource) GetPaymentRequestVotes(c *gin.Context) {
	userId, _ := c.Get("id")
	user := userId.(account.User)

	votes, err := r.dao.GetPaymentRequestVotes(user)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, votes)
}

func (r *DaoResource) UpdatePaymentRequestVotes(c *gin.Context) {
	userId, _ := c.Get("id")
	user := userId.(account.User)

	votesDto := make([]*dto.Vote, 0)
	if err := c.BindJSON(&votesDto); err != nil {
		return
	}

	err := r.dao.UpdatePaymentRequestVotes(votesDto, user)
	if err != nil {
		_ = c.Error(err).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, votesDto)
}
