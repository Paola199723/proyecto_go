package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func listpage(token string, c *gin.Context) {

	url := "https://8j5baasof2.execute-api.us-west-2.amazonaws.com/production/swechallenge/list"
	nextPage := c.Query("next_page")
	if nextPage != "" {
		url = fmt.Sprintf("%s?next_page=%s", url, nextPage)
	}

	reqExternal, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando la solicitud externa"})
		return
	}

	// Configurar encabezados
	token2 := "Bearer " + token
	fmt.Print(token2)
	reqExternal.Header.Set("Authorization", token2)
	reqExternal.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud
	client := &http.Client{}
	resp, err := client.Do(reqExternal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al contactar la API externa"})
		return
	}
	defer resp.Body.Close()

	// Leer la respuesta de la API externa
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al leer respuesta de la API externa"})
		return
	}

	// Mapear el JSON a la estructura AuthResponse
	var authResplist StockResponse

	err2 := json.Unmarshal(body, &authResplist)
	if err2 != nil {
		fmt.Println("Error al parsear JSON:", err)
		return
	}

	if len(authResplist.Items) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No hay datos nuevos para guardar"})
		return
	}
	// Imprimir la lista de acciones
	for _, stock := range authResplist.Items {
		fmt.Printf("Ticker: %s, Empresa: %s, Precio Objetivo: %s -> %s\n",
			stock.Ticker, stock.Company, stock.Target_from, stock.Target_to)
	}

	db := connectDB()

	// Guardar cada acción en la base de datos
	for _, stock := range authResplist.Items {
		db.FirstOrCreate(&stock, ListResponse{Ticker: stock.Ticker, Time: stock.Time})
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "datos guardaddos",
		"total_items": len(authResplist.Items),
		"next_page":   authResplist.NextPage})
}
