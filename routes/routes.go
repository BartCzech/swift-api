package routes

import (
	"github.com/BartCzech/swift-api/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	// Task SWIFT code endpoints
	router.GET("/v1/swift-codes", handlers.GetSwiftCodes)
	router.GET("/v1/swift-codes/:swift-code", handlers.GetSwiftCodeDetails)
	router.GET("/v1/swift-codes/country/:countryISO2code", handlers.GetSwiftCodesByCountry)
	router.POST("/v1/swift-codes", handlers.CreateSwiftCode)
	router.DELETE("/v1/swift-codes/:swift-code", handlers.DeleteSwiftCode)

	return router
}
