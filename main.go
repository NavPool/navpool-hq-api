package main

import (
	"github.com/NavPool/navpool-hq-api/auth"
	"github.com/NavPool/navpool-hq-api/config"
	"github.com/NavPool/navpool-hq-api/database"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	db := database.NewConnection()
	db.AutoMigrate(&auth.User{})

	r := gin.New()

	r.Use(networkSelect)
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to NavPool HQ API!")
	})

	authMiddleware, err := auth.Middleware()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	authGroup := r.Group("/auth")
	authGroup.POST("/login", authMiddleware.LoginHandler)
	authGroup.GET("/refresh_token", authMiddleware.RefreshHandler)

	_ = r.Run(":" + config.Get().Server.Port)
}

func networkSelect(c *gin.Context) {
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
