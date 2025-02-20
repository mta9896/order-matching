package handlers

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(engine *gin.Engine) {
	api := engine.Group("/api") 
	{
		api.POST("/orders", CreateOrder)
		api.GET("/orders", GetOrders)
		api.GET("/orders/:id", GetOrderByID)
	}
}
