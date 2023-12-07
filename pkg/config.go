package pkg

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Host             string
	Port             int64
	PostgresDB       string
	PostgresUser     string
	PostgresPassword string
}

func MustLoad() *Config {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	cfg := Config{
		Host:             getEnv("HOST", "localhost"),
		Port:             getEnvInt("PORT", 5432),
		PostgresDB:       getEnv("POSTGRES_DB", "test"),
		PostgresUser:     getEnv("POSTGRES_USER", "test"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "test"),
	}

	return &cfg
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func getEnvInt(name string, defaultVal int64) int64 {
	valueStr := getEnv(name, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return int64(value)
	}

	return defaultVal
}
