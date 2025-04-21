package database

import (
	"database/sql"
	"fmt"
	"log"

	"gastrobar-backend/config"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB(cfg config.Config) {
	var err error

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error al abrir la conexión con la base de datos: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	log.Println("Conexión a la base de datos exitosa")
}
