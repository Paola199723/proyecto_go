package service

import (
	"log"
	"regexp"
	"sort"
	"strconv"

	"github.com/paola/proyecto_go/configuration"
	"github.com/paola/proyecto_go/internal/models"
)

// Función para limpiar y convertir precios
func CleanAndConvert(priceStr string) float64 {
	re := regexp.MustCompile(`[$,]`) // Remover "$" y ","
	cleanedStr := re.ReplaceAllString(priceStr, "")

	price, err := strconv.ParseFloat(cleanedStr, 64)
	if err != nil {
		log.Printf("Error al convertir %s a float: %v", priceStr, err)
		return -1 // 🔥 Devolver un valor especial en lugar de 0
	}

	return price
}

// Obtener las mejores 3 recomendaciones
func GetTopRecommendations() ([]models.Recommendation, error) {
	var recommendations []models.Recommendation

	query := `SELECT ticker, company, target_from, target_to FROM list_responses`
	db := configuration.GetDB()
	result := db.Raw(query).Scan(&recommendations)
	if result.Error != nil {
		return nil, result.Error
	}

	// Procesar datos
	for i := range recommendations {
		targetFrom := CleanAndConvert(recommendations[i].TargetFrom)
		targetTo := CleanAndConvert(recommendations[i].TargetTo)
		recommendations[i].CambioPorcentual = ((targetTo - targetFrom) / targetFrom) * 100
	}

	// Ordenar y devolver solo los 3 mejores
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].CambioPorcentual > recommendations[j].CambioPorcentual
	})
	if len(recommendations) > 3 {
		recommendations = recommendations[:3]
	}

	return recommendations, nil
}
