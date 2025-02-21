package handlers

import (
	"fmt"
	"net/http"
	"order-matching/models"
	"order-matching/services"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	existingUUIDs  = make(map[string]struct{})
	mutex  sync.Mutex
)

type Response struct {
	Message string `json:"message"`
	Data []models.Order `json:"data"`
}

type OrderBookResponse struct {
	Data []models.OrderBookEntry `json:"data"`
}

func CreateOrder(orderBook *services.OrderBook) gin.HandlerFunc {
	return func(c *gin.Context) {
		var order models.Order
		if err := c.ShouldBindJSON(&order); err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusUnprocessableEntity, Response{
				Message: "Invalid request",
				Data: []models.Order{},
			})
			return
		}

		mutex.Lock()
		defer mutex.Unlock() // handle this mutex in the service too

		if _, exists := existingUUIDs[order.ID]; exists {
			c.JSON(http.StatusConflict, Response{
				Message: "This order has been processed already.",
				Data: []models.Order{},
			})
			return
		}

		existingUUIDs[order.ID] = struct{}{}

		matchedOrders := orderBook.PlaceOrder(&order)
		if matchedOrders == nil {
			matchedOrders = []models.Order{}
		}

		c.JSON(http.StatusOK, Response{
			Message: "success",
			Data: matchedOrders,
		})
	}
}

func GetOrderBook(orderBook *services.OrderBook) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
		if err != nil || limit <= 0 {
			limit = 10
		}

		result := orderBook.GetOrderBook(int(limit))

		c.JSON(http.StatusOK, OrderBookResponse{
			Data: result,
		})
	}
}

func GetOrdersList(orderBook *services.OrderBook) gin.HandlerFunc {
	return func(c *gin.Context) {
		page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
		if err != nil || page <= 0 {
			page = 1
		}

		pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
		if err != nil || pageSize <= 0 {
			pageSize = 10
		}

		orders := orderBook.GetOrderList(page, pageSize)
		fmt.Println("orders", orders)

		c.JSON(http.StatusOK, Response{
			Message: "success",
			Data: orders,
		})
	}
}

func GetOrderByID(c *gin.Context) {

}