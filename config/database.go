package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DataBase *gorm.DB

// Conexión a la base de datos
func connectDB() {
	dsn := "postgresql://paola:BPYlCMSt9BFE7oFW0RpLBw@cuter-example-8394.j77.aws-us-east-1.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}
	fmt.Println("✅ Conectado a CockroachDB")

	// Migrar la estructura a la base de datos (Crea la tabla si no existe)
	db.AutoMigrate(&ListResponse{})
	DataBase = db
}

func getDB() *gorm.DB {

	return DataBase
}
