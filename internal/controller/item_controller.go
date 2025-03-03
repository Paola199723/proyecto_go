package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/paola/proyecto_go/internal/service"
)

func Listpage(c *gin.Context) {
	var request service.ListPageRequest

	// Validar JSON del request
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token no proporcionado o inválido"})
		return
	}

	// Obtener el next_page desde la URL
	nextPage := c.Query("next_page")

	// Llamar al servicio
	authResplist, err := service.FetchAndStoreItems(request, nextPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener datos de la API externa"})
		return
	}

	// Responder al frontend
	c.JSON(http.StatusOK, gin.H{
		"message":     "Datos guardados",
		"total_items": len(authResplist.Items),
		"next_page":   authResplist.NextPage,
	})
}
