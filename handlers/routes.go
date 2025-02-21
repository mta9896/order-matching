package handlers

import (
	"order-matching/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(engine *gin.Engine) {
	orderBook := services.NewOrderBook()
	api := engine.Group("/api") 
	{
		api.POST("/orders", CreateOrder(orderBook))
		api.GET("/orderbook", GetOrderBook(orderBook))
		api.GET("/orders", GetOrdersList(orderBook))
	}
}
