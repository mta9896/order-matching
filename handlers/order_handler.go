package handlers

import (
	"fmt"
	"net/http"
	"order-matching/models"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		fmt.Println("hererere")
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid request"})
		return
	}


}

func GetOrders(c *gin.Context) {

}

func GetOrderByID(c *gin.Context) {

}