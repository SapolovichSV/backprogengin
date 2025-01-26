package main

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SapolovichSV/backprogeng/internal/authmiddleware"
	"github.com/SapolovichSV/backprogeng/internal/config"
	drinkController "github.com/SapolovichSV/backprogeng/internal/drink/controller"
	drinkModel "github.com/SapolovichSV/backprogeng/internal/drink/model"
	httpinfra "github.com/SapolovichSV/backprogeng/internal/http_infra"
	"github.com/SapolovichSV/backprogeng/internal/logger"
	userController "github.com/SapolovichSV/backprogeng/internal/user/controller"
	userModel "github.com/SapolovichSV/backprogeng/internal/user/model"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const SERVER_SHUTDOWN_TIMEOUT = 10 * time.Second

// @title backProgeng API Info
// @version 1.0
// @description This is a simple backend for a out web application
// @BasePath /api
func main() {
	Run()
}
func Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	config := config.ListConfig()
	logger := logger.New(slog.Level(config.LogLevel))
	logger.Info("Config parsed", "config", config)
	//sudo docker run --rm --name db -p 5432:5432 -e POSTGRES_PASSWORD=pass123 -d postgres
	migrateAndUp(&config, logger)

	conn, err := pgxpool.New(ctx, config.DbAddr)
	if err != nil {
		panic("Connection fail" + err.Error())
	}
	defer conn.Close()
	if err := conn.Ping(ctx); err != nil {
		panic("Con't ping db" + err.Error())
	}
	//Создаём модель дринков
	modelDrink := drinkModel.New(conn)
	modelUser := userModel.New(conn)

	//Создаём контроллер дринков
	drinkHandler := drinkController.New(modelDrink, ctx)

	authmiddle := authmiddleware.New()
	userHandler := userController.New(modelUser, authmiddle, ctx)
	//Создаём сервер и в его роутер записываем роуты дринктов и еще юзеров(ещё их не наиписал)
	server := httpinfra.NewServer(config.Port)
	router := server.GetRouter()

	drinkHandler.AddRoutes("api", router)
	userHandler.AddRoutes("api", router)
	//Запускаем сервер
	err = server.Start()
	if err != nil {
		panic(err)
	}
	//также создаём горутину для грасефук шатдауна
	go func() {
		<-ctx.Done()
		logger.Info("Shutting down server")
		ctx, cancel := context.WithTimeout(context.Background(), SERVER_SHUTDOWN_TIMEOUT)
		defer cancel()
		if err := server.Stop(ctx); err != nil {
			logger.Error("Failed to stop server", "error", err)
		}
	}()
}
func migrateAndUp(config *config.Config, logger *slog.Logger) {
	db, err := sql.Open("pgx", config.DbAddr)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	logger.Info("Connected to db")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		panic(err.Error())
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
	defer db.Close()
}
