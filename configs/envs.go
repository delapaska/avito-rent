package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	Host       string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		Port:       getEnv("PORT", "8080"),
		Host:       getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5436"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "avitorent"),
		JWTSecret:  getEnv("JWT_SECRET", "some-secret"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}
