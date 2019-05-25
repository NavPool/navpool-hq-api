package api

import (
	"github.com/NavPool/navpool-hq-api/config"
	"github.com/NavPool/navpool-hq-api/database/migrate"
	"github.com/NavPool/navpool-hq-api/middleware"
	"github.com/NavPool/navpool-hq-api/service/account"
	"github.com/NavPool/navpool-hq-api/service/address"
	"github.com/NavPool/navpool-hq-api/service/auth"
	"github.com/NavPool/navpool-hq-api/service/communityFund"
	"github.com/NavPool/navpool-hq-api/service/network"
	"github.com/NavPool/navpool-hq-api/service/staking"
	"github.com/NavPool/navpool-hq-api/service/twofactor"
	"github.com/getsentry/raven-go"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func main() {
	setReleaseMode()

	migrate.Migrate()

	if config.Get().Sentry.Active {
		_ = raven.SetDSN(config.Get().Sentry.DSN)
	}

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

	networkController := new(network.Controller)
	apiGroup.GET("/network/pool-stats", networkController.GetPoolStats)

	apiGroup.Use(authMiddleware.MiddlewareFunc())
	{
		accountController := new(account.Controller)
		apiGroup.GET("/account", accountController.GetAccount)

		twoFactorController := new(twofactor.Controller)
		apiGroup.GET("/2fa/activate", twoFactorController.GetTwoFactorSecret)
		apiGroup.POST("/2fa/enable", twoFactorController.EnableTwoFactor)
		apiGroup.POST("/2fa/disable", twoFactorController.DisableTwoFactor)

		addressController := new(address.Controller)
		apiGroup.GET("/address/:id", addressController.GetAddress)
		apiGroup.DELETE("/address/:id", addressController.DeleteAddress)
		apiGroup.GET("/address", addressController.GetAddresses)
		apiGroup.POST("/address", addressController.CreateAddress)

		communityFundController := new(communityFund.Controller)
		apiGroup.GET("/community-fund/proposal/vote", communityFundController.GetProposalVotes)
		apiGroup.PUT("/community-fund/proposal/vote", communityFundController.UpdateProposalVotes)
		apiGroup.GET("/community-fund/payment-request/vote", communityFundController.GetPaymentRequestVotes)
		apiGroup.PUT("/community-fund/payment-request/vote", communityFundController.UpdatePaymentRequestVotes)

		stakingController := new(staking.Controller)
		apiGroup.GET("/staking/rewards", stakingController.GetStakingRewardsForAccount)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": 404, "message": "Resource not found"})
	})

	_ = r.Run(":" + config.Get().Server.Port)
}

func setReleaseMode() {
	if config.Get().Debug == false {
		log.Printf("Mode: %s", gin.ReleaseMode)
		gin.SetMode(gin.ReleaseMode)
	} else {
		log.Printf("Mode: %s", gin.DebugMode)
		gin.SetMode(gin.DebugMode)
	}
}
