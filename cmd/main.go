package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// creando el json para el http
type Login struct {
	Username string `json:"username" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"auth_token"`
}

type ListResponse struct {
	ID          uint   `gorm:"primaryKey"`
	Ticker      string `json:"ticker"`
	Target_from string `json:"target_from"`
	Target_to   string `json:"target_to"`
	Company     string `json:"company"`
	Action      string `json:"action"`
	Brokerage   string `json:"brokerage"`
	Rating_from string `json:"rating_from"`
	Rating_to   string `json:"rating_to"`
	Time        string `json:"time"`
}
type StockResponse struct {
	Items    []ListResponse `json:"items"`
	NextPage string         `json:"next_page"`
}

func main() {
	route := gin.Default()
	route.GET("/userlogin", startPage)
	route.GET("/user/recomendations", RecomendationsActions)
	route.Run(":8080")
}

func RecomendationsActions(c *gin.Context) {
	db := connectDB()
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

func startPage(c *gin.Context) {
	var login Login
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
	var authResp AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al parsear respuesta JSON"})
		return
	}

	// Enviar el token recibido al cliente

	listpage(authResp.Token, c)
	//c.JSON(http.StatusOK, gin.H{"auth_token": authResp.Token})
}

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

// Conexión a la base de datos
func connectDB() *gorm.DB {
	dsn := "postgresql://paola:BPYlCMSt9BFE7oFW0RpLBw@cuter-example-8394.j77.aws-us-east-1.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}
	fmt.Println("✅ Conectado a CockroachDB")

	// Migrar la estructura a la base de datos (Crea la tabla si no existe)
	db.AutoMigrate(&ListResponse{})

	return db
}
