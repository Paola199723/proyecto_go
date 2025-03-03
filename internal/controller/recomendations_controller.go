package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paola/proyecto_go/internal/service"
)

func RecomendationsList(c *gin.Context) {
	recommendations, err := service.GetTopRecommendations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener recomendaciones"})
		return
	}

	c.JSON(http.StatusOK, recommendations)

}
