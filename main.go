package main

import (
	"fmt"
	"github.com/NavPool/navpool-hq-api/internal/config"
	"github.com/NavPool/navpool-hq-api/internal/database/migrate"
	"github.com/NavPool/navpool-hq-api/internal/di"
	"github.com/NavPool/navpool-hq-api/internal/framework"
	"github.com/NavPool/navpool-hq-api/internal/resource"
	"github.com/NavPool/navpool-hq-api/internal/service/auth"
	"github.com/getsentry/raven-go"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	di.Init()
	config.Init()

	if config.Get().Debug {
		log.SetLevel(log.DebugLevel)
	}

	framework.SetReleaseMode(config.Get().Debug)

	migrate.Migrate()

	if config.Get().Sentry.Active {
		_ = raven.SetDSN(config.Get().Sentry.DSN)
	}

	r := gin.New()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(framework.Cors())
	r.Use(framework.Options)
	r.Use(framework.ErrorHandler)

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Welcome to NavPool HQ API!")
	})

	authMiddleware, err := auth.Middleware()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
		return
	}

	authResource := new(resource.AuthResource)
	authGroup := r.Group("/auth")
	authGroup.POST("/login", authMiddleware.LoginHandler)
	authGroup.POST("/register", authResource.Register)
	authGroup.GET("/refresh-token", authMiddleware.RefreshHandler)

	apiGroup := r.Group("/")

	apiGroup.Use(authMiddleware.MiddlewareFunc())
	{
		accountResource := resource.NewAccountResource(di.Get().GetAccountService())
		apiGroup.GET("/account", accountResource.GetAccount)

		twoFactorResource := resource.NewTwoFactorResource(di.Get().GetAccountService(), di.Get().GetTwofactorService())
		apiGroup.GET("/2fa/activate", twoFactorResource.GetTwoFactorSecret)
		apiGroup.POST("/2fa/enable", twoFactorResource.EnableTwoFactor)
		apiGroup.POST("/2fa/disable", twoFactorResource.DisableTwoFactor)

		addressResource := resource.NewAddressResource(di.Get().GetAddressService())
		apiGroup.GET("/address/:id", addressResource.GetAddress)
		apiGroup.DELETE("/address/:id", addressResource.DeleteAddress)
		apiGroup.GET("/address", addressResource.GetAddresses)
		apiGroup.POST("/address", addressResource.CreateAddress)

		daoResource := resource.NewDaoResource(di.Get().GetDaoService())
		apiGroup.GET("/community-fund/proposal/vote", daoResource.GetProposalVotes)
		apiGroup.PUT("/community-fund/proposal/vote", daoResource.UpdateProposalVotes)
		apiGroup.GET("/community-fund/payment-request/vote", daoResource.GetPaymentRequestVotes)
		apiGroup.PUT("/community-fund/payment-request/vote", daoResource.UpdatePaymentRequestVotes)
		//Legacy cfund endpoints
		apiGroup.GET("/dao/proposal/vote", daoResource.GetProposalVotes)
		apiGroup.PUT("/dao/proposal/vote", daoResource.UpdateProposalVotes)
		apiGroup.GET("/dao/payment-request/vote", daoResource.GetPaymentRequestVotes)
		apiGroup.PUT("/dao/payment-request/vote", daoResource.UpdatePaymentRequestVotes)

		networkResource := resource.NewNetworkResource(di.Get().GetNetworkService())
		apiGroup.GET("/network/stats", networkResource.GetNetworkStats)

		stakingResource := resource.NewStakingResource(di.Get().GetStakingService())
		apiGroup.GET("/staking/rewards", stakingResource.GetStakingRewardsForAccount)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resource Not Found"})
	})

	_ = r.Run(fmt.Sprintf(":%d", config.Get().Server.Port))
}
