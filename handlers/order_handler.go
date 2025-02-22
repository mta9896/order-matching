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

// CreateOrder places a new order in the order book
//	@Summary		Create a new order
//	@Description	Places a buy or sell order in the order book and returns matched orders if available.
//	@Tags			Orders
//	@Accept			json
//	@Produce		json
//	@Param			order	body		models.Order	true	"Order details"	Example({ "uuid": "550e8400-e29b-41d4-a716-446655440000", "action": "BUY", "price": 100.5, "amount": 2 })
//	@Success		200		{object}	Response		"Order successfully placed"
//	@Failure		422		{object}	Response		"Invalid request payload"
//	@Failure		409		{object}	Response		"Duplicate order detected"
//	@Router			/orders [post]
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
		defer mutex.Unlock()

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

// GetOrderBook retrieves the current state of the order book.
//
//	@Summary		Get order book
//	@Description	Returns a list of buy and sell orders with their price and liquidity.
//	@Tags			Orders
//	@Produce		json
//	@Param			limit	query		int	false	"Number of orders to retrieve (default is 10)"
//	@Success		200		{object}	OrderBookResponse	"Successfully retrieved order book"
//	@Router			/orderbook [get]
//	@Example		{json} Success-Response
//	{
//	  "data": [
//	    {
//	      "price": 100.0,
//	      "liquidity": 5.0,
//	      "type": "BUY"
//	    },
//	    {
//	      "price": 99.5,
//	      "liquidity": 3.0,
//	      "type": "SELL"
//	    }
//	  ]
//	}
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

// GetOrdersList retrieves a paginated list of all orders.
//
//	@Summary		Get list of orders
//	@Description	Returns a paginated list of all orders placed in the order book.
//	@Tags			Orders
//	@Produce		json
//	@Param			page		query	int	false	"Page number (default is 1)"
//	@Param			page_size	query	int	false	"Number of orders per page (default is 10)"
//	@Success		200			{object}	Response	"Successfully retrieved list of orders"
//	@Router			/orders [get]
//	@Example		{json} Success-Response
//	{
//	  "message": "success",
//	  "data": [
//	    {
//	      "uuid": "550e8400-e29b-41d4-a716-446655440000",
//	      "action": "BUY",
//	      "price": 100.0,
//	      "amount": 2.5
//	    },
//	    {
//	      "uuid": "550e8400-e29b-41d4-a716-446655440001",
//	      "action": "SELL",
//	      "price": 99.5,
//	      "amount": 1.0
//	    }
//	  ]
//	}
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