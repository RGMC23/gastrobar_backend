package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

// Config define la estructura para las variables de configuración
type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    Port       string
    JWTSecret  string
    JWTExpire  string
}

var cfg Config // Variable global para almacenar la configuración

// LoadConfig carga las variables de entorno y las almacena en la configuración global
func LoadConfig() Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using default environment variables")
    }

    cfg = Config{
        DBHost:     os.Getenv("DB_HOST"),
        DBPort:     os.Getenv("DB_PORT"),
        DBUser:     os.Getenv("DB_USER"),
        DBPassword: os.Getenv("DB_PASSWORD"),
        DBName:     os.Getenv("DB_NAME"),
        Port:       os.Getenv("PORT"),
        JWTSecret:  os.Getenv("JWT_SECRET"),
        JWTExpire:  os.Getenv("JWT_EXPIRATION"),
    }

    return cfg
}

// GetJWTSecret devuelve el valor de JWTSecret
func GetJWTSecret() string {
    if cfg.JWTSecret == "" {
        log.Fatal("JWT_SECRET not set in configuration")
    }
    return cfg.JWTSecret
}

// GetJWTExpire devuelve el valor de JWTExpire
func GetJWTExpire() string {
    if cfg.JWTExpire == "" {
        log.Fatal("JWT_EXPIRATION not set in configuration")
    }
    return cfg.JWTExpire
}