package handlers

import (
	"bytes"
	"container/heap"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"order-matching/models"
	"order-matching/services"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	t.Run("It returns 422 error for invalid action", func(t *testing.T) {
		t.Parallel()
		body := `{
			"uuid": "550e8400-e29b-41d4-a716-446655440000",
			"action": "invalid_action",
			"price": 10.0,
			"amount": 12.0
		}`

		engine := gin.New()
		engine.POST("/api/orders", CreateOrder(services.NewOrderBook()))

		req, _ := http.NewRequest(http.MethodPost, "/api/orders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		response := new(Response)
		json.Unmarshal(w.Body.Bytes(), response)
		assert.Equal(t, "Invalid request", response.Message)
	})

	t.Run("It returns 422 error if mandatory fields are not provided", func(t *testing.T) {
		t.Parallel()
		body := `{
		}`

		engine := gin.New()
		engine.POST("/api/orders", CreateOrder(services.NewOrderBook()))

		req, _ := http.NewRequest(http.MethodPost, "/api/orders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		
		response := new(Response)
		json.Unmarshal(w.Body.Bytes(), response)
		assert.Equal(t, "Invalid request", response.Message)
	})

	t.Run("It returns 422 error if uuid is not a valid uuid4", func(t *testing.T) {
		t.Parallel()
		body := `{
			"uuid": "550e8400-e29b-41d4-a716-4466554400",
			"action": "BUY",
			"price": 10.0,
			"amount": 12.0
		}`

		engine := gin.New()
		engine.POST("/api/orders", CreateOrder(services.NewOrderBook()))

		req, _ := http.NewRequest(http.MethodPost, "/api/orders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		
		response := new(Response)
		json.Unmarshal(w.Body.Bytes(), response)
		assert.Equal(t, "Invalid request", response.Message)
	})

	t.Run("It returns 409 error if the order is duplicate", func(t *testing.T) {
		t.Parallel()
		body := `{
			"uuid": "550e8400-e29b-41d4-a716-446655441000",
			"action": "BUY",
			"price": 10.0,
			"amount": 12.0
		}`

		engine := gin.New()
		engine.POST("/api/orders", CreateOrder(services.NewOrderBook()))

		req, _ := http.NewRequest(http.MethodPost, "/api/orders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, req)

		//send the same request again
		newReq, _ := http.NewRequest(http.MethodPost, "/api/orders", bytes.NewBufferString(body))
		newReq.Header.Set("Content-Type", "application/json")

		newRecorder := httptest.NewRecorder()
		engine.ServeHTTP(newRecorder, newReq)

		assert.Equal(t, http.StatusConflict, newRecorder.Code)
		
		response := new(Response)
		json.Unmarshal(newRecorder.Body.Bytes(), response)
		assert.Equal(t, "This order has been processed already.", response.Message)
	})

	t.Run("It returns 200 for buy order submission when there are no sell orders", func(t *testing.T) {
		t.Parallel()
		body := `{
			"uuid": "550e8400-e29b-41d4-a716-646655440000",
			"action": "BUY",
			"price": 10.0,
			"amount": 12.0
		}`

		engine := gin.New()
		engine.POST("/api/orders", CreateOrder(services.NewOrderBook()))

		req, _ := http.NewRequest(http.MethodPost, "/api/orders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		response := new(Response)
		err := json.Unmarshal(recorder.Body.Bytes(), response)
		assert.Nil(t, err)
		assert.Equal(t, "success", response.Message)
		assert.Equal(t, []models.Order{}, response.Data)
	})

	t.Run("It returns 200 for sell order submission when there are no buy orders", func(t *testing.T) {
		t.Parallel()
		body := `{
			"uuid": "550e8400-e29b-41d4-a716-646655440300",
			"action": "SELL",
			"price": 10.0,
			"amount": 12.0
		}`

		engine := gin.New()
		engine.POST("/api/orders", CreateOrder(services.NewOrderBook()))

		req, _ := http.NewRequest(http.MethodPost, "/api/orders", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)

		response := new(Response)
		err := json.Unmarshal(recorder.Body.Bytes(), response)
		assert.Nil(t, err)
		assert.Equal(t, "success", response.Message)
		assert.Equal(t, []models.Order{}, response.Data)
	})
}

func TestOrderBook(t *testing.T) {
	t.Parallel()
	t.Run("It returns orderbook correctly", func(t *testing.T) {
		t.Parallel()
		orderBook := initOrderBook()

		engine := gin.New()
		engine.GET("/api/orderbook", GetOrderBook(orderBook))

		req, _ := http.NewRequest(http.MethodGet, "/api/orderbook", nil)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		
		response := new(Response)
		json.Unmarshal(recorder.Body.Bytes(), response)
		assert.Equal(t, 4, len(response.Data))
	})
}

func TestOrderList(t *testing.T) {
	t.Parallel()
	t.Run("It returns order list correctly", func(t *testing.T) {
		t.Parallel()
		orderBook := initOrderBook()

		engine := gin.New()
		engine.GET("/api/orders", GetOrdersList(orderBook))

		req, _ := http.NewRequest(http.MethodGet, "/api/orders", nil)
		req.Header.Set("Content-Type", "application/json")

		recorder := httptest.NewRecorder()
		engine.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
		
		response := new(Response)
		json.Unmarshal(recorder.Body.Bytes(), response)
		assert.Equal(t, 6, len(response.Data))
	})
}

func initOrderBook() *services.OrderBook {
	orderBook := services.NewOrderBook()
	orderBook.SellOrders = map[float64][]models.Order{
		120.0: {
			{
				ID:     "550e8400-e29b-41d4-a716-666655442000",
				Action: models.Sell,
				Price:  120.0,
				Amount: 2.0,
			},
		},
		100.0: {
			{
				ID:     "550e8400-e29b-41d4-a716-77755442000",
				Action: models.Sell,
				Price:  100.0,
				Amount: 2.0,
			},
			{
				ID:     "550e8400-e29b-41d4-a716-77755442001",
				Action: models.Sell,
				Price:  100.0,
				Amount: 2.0,
			},
			{
				ID:     "550e8400-e29b-41d4-a716-77755442002",
				Action: models.Sell,
				Price:  100.0,
				Amount: 3.0,
			},
		},
	}

	heap.Push(&orderBook.SellPricesHeap, 120.0)
	heap.Push(&orderBook.SellPricesHeap, 100.0)

	orderBook.BuyOrders = map[float64][]models.Order{
		80.0: {
			{
				ID:     "550e8400-e29b-41d4-a716-666655442000",
				Action: models.Buy,
				Price:  80.0,
				Amount: 2.0,
			},
		},
		100.0: {
			{
				ID:     "550e8400-e29b-41d4-a716-77755442000",
				Action: models.Buy,
				Price:  100.0,
				Amount: 2.0,
			},
		},
	}

	heap.Push(&orderBook.BuyPricesHeap, 80.0)
	heap.Push(&orderBook.BuyPricesHeap, 100.0)

	return orderBook
}