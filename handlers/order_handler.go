package handlers

import (
	"fmt"
	"net/http"
	"order-matching/models"
	"order-matching/services"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	orderBook  = services.NewOrderBook()
	existingUUIDs  = make(map[string]struct{})
	mutex  sync.Mutex
)

type Response struct {
	Message string `json:"message"`
	Data []models.Order `json:"data"`
}

func CreateOrder(c *gin.Context) {
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

func GetOrders(c *gin.Context) {

}

func GetOrderByID(c *gin.Context) {

}