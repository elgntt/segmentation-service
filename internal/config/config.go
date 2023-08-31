package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	PgUser     string
	PgPassword string
	PgHost     string
	PgPort     uint16
	PgDatabase string
	PgSSLMode  string
}

type ServerConfig struct {
	HTTPPort       string
	ServerEndpoint string
}

func GetDBConfig() (DBConfig, error) {
	pgPort, err := strconv.ParseInt(getKey("PGPORT"), 0, 16)
	if err != nil {
		return DBConfig{}, err
	}

	return DBConfig{
		PgUser:     getKey("PGUSER"),
		PgPassword: getKey("PGPASSWORD"),
		PgHost:     getKey("PGHOST"),
		PgPort:     uint16(pgPort),
		PgDatabase: getKey("PGDATABASE"),
		PgSSLMode:  getKey("PGSSLMODE"),
	}, nil
}

func GetServerConfig() ServerConfig {
	return ServerConfig{
		HTTPPort:       ":" + getKey("HTTP_PORT"),
		ServerEndpoint: getKey("SERVER_ENDPOINT"),
	}
}

func getKey(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Error loading .env file")
		return ""
	}

	return os.Getenv(key)
}
