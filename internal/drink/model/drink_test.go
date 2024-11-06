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

func TestSQLDrinkModel_UpdateDrink(t *testing.T) {
	db, err := sql.Open("pgx", "host=localhost user=username password=password dbname=dbname sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()
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

	drink.Tags = []string{"tag3", "tag4"}

	_, err = model.UpdateDrink(ctx, drink)
	if err != nil {
		t.Fatalf("Failed to update drink: %v", err)
	}

	var updatedDrink Drink

	err = db.QueryRowContext(ctx, "SELECT name, tags FROM drinks WHERE name = $1", drink.Name).Scan(&updatedDrink.name, &updatedDrink.tags)
	if err != nil {
		t.Fatalf("Failed to retrieve updated drink: %v", err)
	}
	resDrink := fromModelToController(updatedDrink)
	if !reflect.DeepEqual(drink, resDrink) {
		t.Errorf("Updated drink does not match: got %v, want %v", resDrink, drink)
	}
}

func TestSQLDrinkModel_DeleteDrink(t *testing.T) {
	db, err := sql.Open("pgx", "host=localhost user=username password=password dbname=dbname sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()
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

	err = model.DeleteDrink(ctx, drink.Name)
	if err != nil {
		t.Fatalf("Failed to delete drink: %v", err)
	}

	var deletedDrink Drink

	err = db.QueryRowContext(ctx, "SELECT name, tags FROM drinks WHERE name = $1", drink.Name).Scan(&deletedDrink.name, &deletedDrink.tags)
	if err == nil {
		t.Fatalf("Deleted drink was found: %v", err)
	}
	err = model.DeleteDrink(ctx, drink.Name)
	if err != ErrNotFound {
		t.Fatalf("Error introspection failed: %v", err)
	}
}

func TestSQLDrinkModel_DrinksByTags(t *testing.T) {
	db, err := sql.Open("pgx", "host=localhost user=username password=password dbname=dbname sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()
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
	drink1 := entities.Drink{
		Name: "Test Drink1",
		Tags: []string{"tag1", "tag2"},
	}
	drink2 := entities.Drink{
		Name: "Test Drink2",
		Tags: []string{"tag2", "tag3"},
	}

	_, err = model.CreateDrink(ctx, drink1)
	if err != nil {
		t.Fatalf("Failed to create drink1: %v", err)
	}
	_, err = model.CreateDrink(ctx, drink2)
	if err != nil {
		t.Fatalf("Failed to create drink2: %v", err)
	}

	drinks, err := model.DrinksByTags(ctx, []string{"tag2"})
	if err != nil {
		t.Fatalf("Failed to get drinks by tags: %v", err)
	}
	if len(drinks) != 2 {
		t.Fatalf("Wrong number of drinks: got %d, want %d", len(drinks), 2)
	}
	if !reflect.DeepEqual(drinks[0], drink1) || !reflect.DeepEqual(drinks[1], drink2) {
		t.Errorf("Drinks do not match: got %v, want %v", drinks, []entities.Drink{drink1, drink2})
	}
}

func TestSQLDrinkModel_AllDrinks(t *testing.T) {
	db, err := sql.Open("pgx", "host=localhost user=username password=password dbname=dbname sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()
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
	drink1 := entities.Drink{
		Name: "Test Drink1",
		Tags: []string{"tag1", "tag2"},
	}
	drink2 := entities.Drink{
		Name: "Test Drink2",
		Tags: []string{"tag2", "tag3"},
	}

	_, err = model.CreateDrink(ctx, drink1)
	if err != nil {
		t.Fatalf("Failed to create drink1: %v", err)
	}
	_, err = model.CreateDrink(ctx, drink2)
	if err != nil {
		t.Fatalf("Failed to create drink2: %v", err)
	}

	drinks, err := model.AllDrinks(ctx, 0)
	if err != nil {
		t.Fatalf("Failed to get all drinks: %v", err)
	}
	if len(drinks) != 2 {
		t.Fatalf("Wrong number of drinks: got %d, want %d", len(drinks), 2)
	}
	if !reflect.DeepEqual(drinks[0], drink1) || !reflect.DeepEqual(drinks[1], drink2) {
		t.Errorf("Drinks do not match: got %v, want %v", drinks, []entities.Drink{drink1, drink2})
	}
}
