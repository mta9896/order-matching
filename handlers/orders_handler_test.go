package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	t.Run("It returns 422 error for invalid action", func(t *testing.T) {
		body := `{
			"uuid": "5f19fe3f-8263-4f7f-8d2e",
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

		fmt.Println(w.Body.String())
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.JSONEq(t, `{"error":"Invalid request"}`, w.Body.String())
	})
	

}