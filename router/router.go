package router

import (
	"net/http"

	"github.com/HiN3i1/wallet-service/controller"
	v1API "github.com/HiN3i1/wallet-service/controller/v1"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// New router
func New() *gin.Engine {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowCredentials = true
	config.AllowOriginFunc = func(origin string) bool {
		return true
	}

	router.NoMethod(controller.HandleNotFound)
	router.NoRoute(controller.HandleNotFound)
	router.Use(cors.New(config))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "server is running"})
	})

	APIv1 := router.Group("/api/v1/")
	{
		walletGroup := APIv1.Group("/wallets")
		{
			walletGroup.POST("/:wallet_id", v1API.SetAPIToken)
			walletGroup.POST("/callback", v1API.Callback)
		}
		customerGroup := APIv1.Group("/customer")
		{
			customerGroup.POST("/apitoken/:coin", v1API.SetAPIToken)
			customerGroup.GET("/:customer_id/:coin", v1API.GetDepositWalletAddresses)
			customerGroup.POST("/:customer_id/:coin", v1API.CreateDepositWalletAddresses)
		}
	}
	return router
}
