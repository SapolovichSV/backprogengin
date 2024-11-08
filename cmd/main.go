package main

import (
	"context"
	"database/sql"

	"github.com/SapolovichSV/backprogeng/internal/config"
	"github.com/SapolovichSV/backprogeng/internal/drink/controller"
	"github.com/SapolovichSV/backprogeng/internal/drink/model"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	Run()
}
func Run() {
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
	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		panic(err.Error())
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		panic(err)
	}
	defer db.Close()
	model := model.NewSQLDrinkModel(db)
	ctx := context.Background()
	ctr := controller.NewHTTPHandler(model, ctx)
	router := ctr.BuildRouter("/api")
	ctr.AddRoutes(router)
	err = ctr.Start(config.Port)
	if err != nil {
		panic("error starting server")
	}
}
