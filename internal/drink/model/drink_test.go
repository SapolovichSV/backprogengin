package model

import (
	"context"
	"database/sql"
	"reflect"
	"testing"

	"github.com/SapolovichSV/backprogeng/internal/drink/entities"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Сначала нужно поднять тестовую бд
// потом запускать
// sudo docker run --name --rm test-postgres -e POSTGRES_PASSWORD=password -e POSTGRES_USER=username -e POSTGRES_DB=dbname -p 5432:5432 -d postgres
func TestSQLDrinkModel_CreateDrink(t *testing.T) {
	db, err := sql.Open("pgx", "host=localhost user=username password=password dbname=dbname sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()
	_, err = db.Exec("DROP TABLE IF EXISTS drinks;")
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}
	_, err = db.Exec(`CREATE TABLE drinks (
	    id SERIAL PRIMARY KEY,
	    name VARCHAR(255) NOT NULL,
	    tags TEXT
	);`)
	defer db.Exec("DROP TABLE IF EXISTS drinks;")
	if err != nil {
		t.Fatalf("Failed to create table: %v", err)
	}
	model := &SQLDrinkModel{db: db}

	ctx := context.Background()
	drink := entities.Drink{
		Name: "Test Drink",
		Tags: []string{"tag1", "tag2"},
	}

	_, err = model.CreateDrink(ctx, drink)
	if err != nil {
		t.Fatalf("Failed to create drink: %v", err)
	}

	var createdDrink Drink

	err = db.QueryRowContext(ctx, "SELECT name, tags FROM drinks WHERE name = $1", drink.Name).Scan(&createdDrink.name, &createdDrink.tags)
	if err != nil {
		t.Fatalf("Failed to retrieve created drink: %v", err)
	}
	resDrink := fromModelToController(createdDrink)
	if !reflect.DeepEqual(drink, resDrink) {
		t.Errorf("Created drink does not match: got %v, want %v", resDrink, drink)
	}
}
