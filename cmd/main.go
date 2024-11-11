package main

import (
	"context"
	"database/sql"

	"github.com/SapolovichSV/backprogeng/internal/config"
	"github.com/SapolovichSV/backprogeng/internal/drink/controller"
	"github.com/SapolovichSV/backprogeng/internal/drink/model"
	httpinfra "github.com/SapolovichSV/backprogeng/internal/http_infra"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	Run()
}
func Run() {
	//Чекаем переменные окружения чтобы подцепить бд и порт сервера
	config := config.ListConfig()
	//sudo docker run --rm --name db -p 5432:5432 -e POSTGRES_PASSWORD=pass123 -d postgres
	db, err := sql.Open("pgx", config.DbAddr)
	if err != nil {
		panic(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}
	//Мигрируем на новую схему бдшки(надо прочекать будет теряются ли данные при этом)
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		panic(err.Error())
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
	defer db.Close()
	//Создаём модель дринков
	model := model.NewSQLDrinkModel(db)
	//Создаём контроллер дринков
	drinkHandler := controller.NewHTTPHandler(model, context.Background())
	//Создаём сервер и в его роутер записываем роуты дринктов и еще юзеров(ещё их не наиписал)
	server := httpinfra.NewServer(config.Port)
	router := server.GetRouter()
	drinkHandler.AddRoutes("api", router)
	//userHandler.AddRoutes("api", router)
	//GameHandler.AddRoutes("api",router))
	//Запускаем сервер
	err = server.Start()
	if err != nil {
		panic(err)
	}
}
