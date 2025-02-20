package handlers

import (
	"fmt"
	"net/http"
	"order-matching/models"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	orderBook  = make(map[string]models.Order)
	existingUUIDs  = make(map[string]struct{})
	mutex  sync.Mutex
)

func CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request"})
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := existingUUIDs[order.ID]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "This order has been processed already."})
	}

	existingUUIDs[order.ID] = struct{}{}

	

}

func GetOrders(c *gin.Context) {

}

func GetOrderByID(c *gin.Context) {

}