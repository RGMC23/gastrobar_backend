package main

import (
	"log"
	"net/http"

	"gastrobar-backend/api"
	"gastrobar-backend/cmd/app"
	"gastrobar-backend/config"
	database "gastrobar-backend/db"
)

func main() {
	// Cargar la configuración desde el archivo .env
	cfg := config.LoadConfig()

	// Conectar a la base de datos
	database.ConnectDB(cfg)

	//Configurar la aplicación
	app := app.NewApp(database.DB)

	// Configurar las rutas
	router := api.SetupRoutes(app)

	// Iniciar el servidor
	log.Printf("Backend running on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, router))
}
