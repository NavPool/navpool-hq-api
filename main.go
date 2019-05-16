package main

import (
	"github.com/NavPool/navpool-hq-api/config"
	"github.com/NavPool/navpool-hq-api/database"
	"github.com/NavPool/navpool-hq-api/middleware"
	"github.com/NavPool/navpool-hq-api/service/account"
	"github.com/NavPool/navpool-hq-api/service/address"
	"github.com/NavPool/navpool-hq-api/service/address/model"
	"github.com/NavPool/navpool-hq-api/service/auth"
	"github.com/NavPool/navpool-hq-api/service/communityFund"
	model2 "github.com/NavPool/navpool-hq-api/service/communityFund/model"
	"github.com/NavPool/navpool-hq-api/service/twofactor"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	if config.Get().Debug == false {
		gin.SetMode(gin.ReleaseMode)
	}

	dbFixtures()

	r := gin.New()

	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.Use(middleware.NetworkSelect)
	r.Use(middleware.Options)
	r.Use(middleware.ErrorHandler)

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to NavPool HQ API!")
	})

	authMiddleware, err := auth.Middleware()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	authController := new(auth.Controller)
	authGroup := r.Group("/auth")
	authGroup.POST("/login", authMiddleware.LoginHandler)
	authGroup.POST("/register", authController.Register)
	authGroup.GET("/refresh-token", authMiddleware.RefreshHandler)

	apiGroup := r.Group("/")
	apiGroup.Use(authMiddleware.MiddlewareFunc())
	{
		authController := new(account.Controller)
		apiGroup.GET("/account", authController.GetAccount)

		twofactorController := new(twofactor.Controller)
		apiGroup.GET("/2fa/activate", twofactorController.GetTwoFactorSecret)
		apiGroup.POST("/2fa/enable", twofactorController.EnableTwoFactor)
		apiGroup.POST("/2fa/disable", twofactorController.DisableTwoFactor)

		addressController := new(address.Controller)
		apiGroup.DELETE("/address/:id", addressController.DeleteAddress)
		apiGroup.GET("/address", addressController.GetAddresses)
		apiGroup.POST("/address", addressController.CreateAddress)

		communityFundController := new(communityFund.Controller)
		apiGroup.GET("/community-fund/proposal/vote", communityFundController.GetProposalVotes)
		apiGroup.PUT("/community-fund/proposal/vote", communityFundController.UpdateProposalVotes)
		apiGroup.GET("/community-fund/payment-request/vote", communityFundController.GetPaymentRequestVotes)
		apiGroup.PUT("/community-fund/payment-request/vote", communityFundController.UpdatePaymentRequestVotes)
	}

	_ = r.Run(":" + config.Get().Server.Port)
}

func dbFixtures() {
	db, err := database.NewConnection()
	if err != nil {
		return
	}

	db.AutoMigrate(&account.User{}, &account.TwoFactor{}, model.Address{}, model2.Vote{})

	if config.Get().Fixtures == true {
		account.CreateUser("admin", "admin")
		account.CreateUser("deleted", "deleted")
		account.DeleteUser("deleted", true)
	}
}
