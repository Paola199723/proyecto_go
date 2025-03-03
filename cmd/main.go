package main

import (
	"github.com/gin-gonic/gin"
	"github.com/paola/proyecto_go/configuration"
	"github.com/paola/proyecto_go/internal/controller"
)

func main() {
	configuration.ConnectDB()
	route := gin.Default()

	route.POST("/user/login", controller.StartPage)
	route.GET("/user/stocks", controller.Listpage)
	route.GET("/user/recomendations", controller.RecomendationsList)
	route.Run("0.0.0.0:8080")
	//route.Run(":8080")
}
