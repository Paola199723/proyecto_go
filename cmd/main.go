package main

import (
	"github.com/gin-gonic/gin"
	"github.com/paola/proyecto_go/internal/controller"
)

func main() {
	route := gin.Default()

	route.GET("/userlogin", controller.StartPage)
	route.GET("/user/recomendations", controller.RecomendationsActions)
	route.Run(":8080")
}
