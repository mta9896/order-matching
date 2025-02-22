package main

import (
	"fmt"
	"net/http"
	"order-matching/handlers"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files" 
	_ "order-matching/docs"
)

//	@title			Order Matching API
//	@version		1.0
//	@description	This is an API for an order matching system.
//	@host			localhost:8080
//	@BasePath		/api
//  @schemes		http

func main() {
	engine := gin.New()
	handlers.RegisterRoutes(engine)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	fmt.Println("Server started on port 8080")
    http.ListenAndServe(":8080", engine)
}
