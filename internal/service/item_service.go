package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/paola/proyecto_go/configuration"
	"github.com/paola/proyecto_go/internal/models"
)

// Definir estructura para recibir el token

// Función para obtener y guardar datos de la API externa
func FetchAndStoreItems(Token string, nextPage string) (models.StockResponse, error) {
	url := "https://8j5baasof2.execute-api.us-west-2.amazonaws.com/production/swechallenge/list"
	if nextPage != "" {
		url = fmt.Sprintf("%s?next_page=%s", url, nextPage)
	}

	// Crear la solicitud HTTP
	reqExternal, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return models.StockResponse{}, err
	}

	// Configurar encabezados
	reqExternal.Header.Set("Authorization", Token)
	println(Token)
	reqExternal.Header.Set("Content-Type", "application/json")

	// Enviar la solicitud
	client := &http.Client{}
	resp, err := client.Do(reqExternal)
	if err != nil {
		return models.StockResponse{}, err
	}
	defer resp.Body.Close()

	// Leer la respuesta de la API externa
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.StockResponse{}, err
	}

	// Depuración: Imprimir la respuesta JSON recibida antes de deserializar
	println("Respuesta JSON recibida:", string(body))

	// Intentar deserializar
	var authResplist models.StockResponse
	if err := json.Unmarshal(body, &authResplist); err != nil {
		println("Error al deserializar JSON:", err.Error())
		return models.StockResponse{}, err
	}

	// Depuración: Verificar si la lista de items está vacía
	if len(authResplist.Items) == 0 {
		println("⚠ Advertencia: La API devolvió una lista vacía")
	} else {
		println("✅ Se recibieron", len(authResplist.Items), "elementos")
	}

	// Leer la respuesta de la API externa
	/*body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.StockResponse{}, err
	}

	// Mapear el JSON a la estructura StockResponse
	var authResplist models.StockResponse

	if err := json.Unmarshal(body, &authResplist); err != nil {
		return models.StockResponse{}, err
	}*/

	// Guardar en la base de datos
	db := configuration.GetDB()

	if len(authResplist.Items) > 0 {
		for _, stock := range authResplist.Items {
			db.FirstOrCreate(&stock, models.ListResponse{Ticker: stock.Ticker, Time: stock.Time})
		}
		// 🔥 No intenta guardar si la lista está vacía
	}

	return authResplist, nil
}
