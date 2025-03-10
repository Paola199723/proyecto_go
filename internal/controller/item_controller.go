package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paola/proyecto_go/internal/service"
)

func Listpage(c *gin.Context) {
	//var request service.ListPageRequest
	authToken := c.GetHeader("Authorization")
	if authToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token no proporcionado"})
		return
	}

	// Obtener el next_page desde la URL
	nextPage := c.Query("next_page")
	//request := service.ListPageRequest{Token: authToken}

	// Llamar al servicio
	authResplist, err := service.FetchAndStoreItems(authToken, nextPage)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener datos de la API externa"})
		return
	}

	// Responder al frontend
	c.JSON(http.StatusOK, gin.H{
		"message":     "Datos guardados",
		"total_items": authResplist.Items,
		"next_page":   authResplist.NextPage,
	})
}
