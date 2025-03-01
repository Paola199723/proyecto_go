package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	models "github.com/paola/proyecto_go/internal/models"
)

func StartPage(c *gin.Context) {
	var login models.Login
	// If `GET`, only `Form` binding engine (`query`) used.
	// If `POST`, first checks the `content-type` for `JSON` or `XML`, then uses `Form` (`form-data`).
	// See more at https://github.com/gin-gonic/gin/blob/master/binding/binding.go#L48
	fmt.Printf("%s , %s", login.Username, login.Password)

	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos", "datos": login})
		return
	}

	// Convertir la solicitud en JSON
	jsonData, err := json.Marshal(login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al procesar JSON"})
		return
	}

	url := "https://8j5baasof2.execute-api.us-west-2.amazonaws.com/production/swechallenge/login"

	reqExternal, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando la solicitud externa"})
		return
	}

	// Configurar encabezados
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
	var authResp models.AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al parsear respuesta JSON"})
		return
	}

	// Enviar el token recibido al cliente

	//listpage(authResp.Token, c)
	c.JSON(http.StatusOK, gin.H{"auth_token": authResp.Token})
}
