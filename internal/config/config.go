package config

import (
	"os"
)

type Config struct {
	Port   string
	DbAddr string
}

func ListConfig() Config {
	port := os.Getenv("PORT")
	dbAddr := os.Getenv("DB_ADDR")
	if port == "" {
		port = "8080"
	}
	if dbAddr == "" {
		dbAddr = "host=db port=5432 user=postgres password=pass123 dbname=postgres sslmode=disable"
	}
	return Config{
		Port:   port,
		DbAddr: dbAddr,
	}
}
