package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	route := gin.Default()
	route.GET("/userlogin", startPage)
	route.GET("/user/recomendations", RecomendationsActions)
	route.Run(":8080")
}

func RecomendationsActions(c *gin.Context) {
	db := getDB()
	var recommendations []struct {
		Ticker           string  `json:"ticker"`
		Company          string  `json:"company"`
		TargetFrom       float64 `json:"target_from"`
		TargetTo         float64 `json:"target_to"`
		CambioPorcentual float64 `json:"cambio_porcentual"`
	}

	// Consulta SQL con CockroachDB
	query := `
        SELECT ticker, company, target_from, target_to,
               ((target_to - target_from) / target_from) * 100 AS cambio_porcentual
        FROM list_responnses
        WHERE target_to > target_from
        ORDER BY cambio_porcentual DESC
        LIMIT 5;
    `

	// Ejecutar la consulta
	result := db.Raw(query).Scan(&recommendations)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al obtener recomendaciones"})
		return
	}

	// Verificar si hay resultados
	if len(recommendations) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No hay recomendaciones disponibles"})
		return
	}

	// Devolver recomendaciones
	c.JSON(http.StatusOK, recommendations)

}
