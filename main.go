package main

import (
	"fmt"
	"net/http"
	"order-matching/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.New()
	handlers.RegisterRoutes(engine)
	fmt.Println("Server started on port 8080")
    http.ListenAndServe(":8080", engine)
}
