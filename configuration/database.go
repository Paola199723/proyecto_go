package configuration

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DataBase *gorm.DB

// Conexión a la base de datos
func ConnectDB() {
	dsn := "postgresql://paola:BPYlCMSt9BFE7oFW0RpLBw@cuter-example-8394.j77.aws-us-east-1.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full"
	var err error
	DataBase, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}) // 🔥 Se asigna correctamente

	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}

	fmt.Println("✅ Conectado a CockroachDB")
	// Migrar la estructura a la base de datos (Crea la tabla si no existe)
	//DataBase.AutoMigrate(&models.ListResponse{})

}

func GetDB() *gorm.DB {
	if DataBase == nil {
		log.Fatal("❌ Error: La base de datos no está inicializada. ¿Llamaste ConnectDB() en main.go?")
	}
	return DataBase
}
