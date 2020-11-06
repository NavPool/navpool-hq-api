package resource

import (
	"github.com/NavPool/navpool-hq-api/internal/service"
	"github.com/gin-gonic/gin"
)

type NetworkResource struct {
	network *service.NetworkService
}

func NewNetworkResource(networkService *service.NetworkService) *NetworkResource {
	return &NetworkResource{networkService}
}

func (r *NetworkResource) GetNetworkStats(c *gin.Context) {
	poolStats, err := r.network.GetNetworkStats()
	if err != nil {
		_ = c.Error(service.ErrorUnableToRetrieveStats).SetType(gin.ErrorTypePublic)
		return
	}

	c.JSON(200, poolStats)
}
