package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"order-matching/models"
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
		engine.POST("/api/orders", CreateOrder)

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
		engine.POST("/api/orders", CreateOrder)

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
		engine.POST("/api/orders", CreateOrder)

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
			"uuid": "550e8400-e29b-41d4-a716-446655440000",
			"action": "BUY",
			"price": 10.0,
			"amount": 12.0
		}`

		engine := gin.New()
		engine.POST("/api/orders", CreateOrder)

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
			"uuid": "550e8400-e29b-41d4-a716-446655440000",
			"action": "BUY",
			"price": 10.0,
			"amount": 12.0
		}`

		engine := gin.New()
		engine.POST("/api/orders", CreateOrder)

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