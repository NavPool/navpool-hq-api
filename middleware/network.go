package middleware

import (
	"github.com/NavPool/navpool-hq-api/config"
	"github.com/gin-gonic/gin"
)

func NetworkSelect(c *gin.Context) {
	switch network := c.GetHeader("Network"); network {
	case "testnet":
		config.Get().SelectedNetwork = network
		break
	case "mainnet":
		config.Get().SelectedNetwork = network
		break
	default:
		config.Get().SelectedNetwork = "mainnet"
	}

	c.Header("X-Network", config.Get().SelectedNetwork)
}
