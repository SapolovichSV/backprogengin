package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port     string
	DbAddr   string
	LogLevel int
}

func ListConfig() Config {
	port := os.Getenv("PORT")

	logLevel := os.Getenv("LOG_LEVEL")
	if port == "" {
		port = "8080"
	}
	// if dbAddr == "" {
	// 	dbAddr = "host=db port=5432 user=postgres password=pass123 dbname=postgres sslmode=disable"
	// }
	if logLevel == "" {
		logLevel = "-4"
	}
	//url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
	//  pql.conf.Postgres.User,
	//    pql.conf.Postgres.Password,
	//    pql.conf.Postgres.Host,
	//    pql.conf.Postgres.Port,
	//    pql.conf.Postgres.DB)
	logLevelInt, err := strconv.Atoi(logLevel)
	incorrectLevel := false
	switch logLevelInt {
	case -4:
		break
	case 0:
		break
	case 4:
		break
	case 8:
		break
	default:
		incorrectLevel = true
	}
	if err != nil || incorrectLevel {
		panic("Incorrect log level from env")
	}
	dbAddr := parseDbAddr()
	return Config{
		Port:     port,
		DbAddr:   dbAddr,
		LogLevel: logLevelInt,
	}
}
func parseDbAddr() string {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres"
	}
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "pass123"
	}
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "postgres"
	}
	dbAddr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
	)
	return dbAddr
}
