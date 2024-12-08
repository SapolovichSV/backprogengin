package main

import (
	"context"
	"database/sql"

	"github.com/SapolovichSV/backprogeng/internal/config"
	drinkController "github.com/SapolovichSV/backprogeng/internal/drink/controller"
	drinkModel "github.com/SapolovichSV/backprogeng/internal/drink/model"
	httpinfra "github.com/SapolovichSV/backprogeng/internal/http_infra"
	userController "github.com/SapolovichSV/backprogeng/internal/user/controller"
	userModel "github.com/SapolovichSV/backprogeng/internal/user/model"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	// echo-swagger middleware
)

// @title Swagger Example
func main() {
	Run()
}
func Run() {
	//Чекаем переменные окружения чтобы подцепить бд и порт сервера
	config := config.ListConfig()
	//sudo docker run --rm --name db -p 5432:5432 -e POSTGRES_PASSWORD=pass123 -d postgres
	migrateAndUp(&config)

	ctx := context.Background()
	conn, err := pgxpool.New(ctx, config.DbAddr)
	if err := conn.Ping(ctx); err != nil {
		panic("Connection fail" + err.Error())
	}
	//Создаём модель дринков
	modelDrink := drinkModel.New(conn)
	modelUser := userModel.New(conn)

	//Создаём контроллер дринков
	drinkHandler := drinkController.New(modelDrink, ctx)
	userHandler := userController.New(modelUser, ctx)
	//Создаём сервер и в его роутер записываем роуты дринктов и еще юзеров(ещё их не наиписал)
	server := httpinfra.NewServer(config.Port)
	router := server.GetRouter()

	drinkHandler.AddRoutes("api", router)
	userHandler.AddRoutes("api", router)
	//userHandler.AddRoutes("api", router)
	//GameHandler.AddRoutes("api",router))
	//Запускаем сервер
	err = server.Start()
	if err != nil {
		panic(err)
	}
}
func migrateAndUp(config *config.Config) {
	db, err := sql.Open("pgx", config.DbAddr)
	if err != nil {
		panic(err)
	}
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
