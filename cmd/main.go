package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/SapolovichSV/backprogeng/internal/config"
	"github.com/SapolovichSV/backprogeng/internal/drink/controller"
	"github.com/SapolovichSV/backprogeng/internal/drink/model"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	config := config.ListConfig()
	db, err := sql.Open("pgx", config.DbAddr)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	model := model.NewSQLDrinkModel(db)
	ctx := context.TODO()
	ctr := controller.NewHTTPHandler(model, ctx)
	router := ctr.BuildRouter("/api")
	ctr.AddRoutes(router)
	go func() {
		ctr.Start(config.Port)
		if err != nil {
			panic("error starting server")
		}
	}()
	time.Sleep(10000 * time.Minute)
}
