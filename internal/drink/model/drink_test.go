package model

import (
	"context"
	"reflect"
	"testing"

	"github.com/SapolovichSV/backprogeng/internal/drink/entities"
	"github.com/jackc/pgx/v5/pgxpool"
)

const QUERY_CREATE_TABLES = `CREATE TABLE drinks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    tags TEXT
);
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255)
);
CREATE TABLE favs (
    user_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)  ON DELETE CASCADE,
    drink_id INT NOT NULL,
   FOREIGN KEY (drink_id) REFERENCES drinks(id) ON DELETE CASCADE
);`
const QUERY_DROP_TABLES = `DROP TABLE drinks CASCADE;
DROP TABLE users CASCADE;
DROP TABLE favs CASCADE;`

// Сначала нужно поднять тестовую бд
// потом запускать
// sudo docker run --rm --name test-postgres -e POSTGRES_PASSWORD=password -e POSTGRES_USER=username -e POSTGRES_DB=dbname -p 5432:5432 -d postgres
func TestSQLDrinkModel_CreateDrink(t *testing.T) {
	db, err := pgxpool.New(context.TODO(), "host=localhost user=username password=password dbname=dbname sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}
	_, err = db.Exec(context.TODO(), QUERY_CREATE_TABLES)
	defer func() {
		_, err := db.Exec(context.TODO(), QUERY_DROP_TABLES)
		if err != nil {
			panic(err)
		}
	}()
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

	err = db.QueryRow(ctx, "SELECT name, tags FROM drinks WHERE name = $1", drink.Name).Scan(&createdDrink.name, &createdDrink.tags)
	if err != nil {
		t.Fatalf("Failed to retrieve created drink: %v", err)
	}
	resDrink := fromModelToController(createdDrink)
	if !reflect.DeepEqual(drink, resDrink) {
		t.Errorf("Created drink does not match: got %v, want %v", resDrink, drink)
	}
}

func TestSQLDrinkModel_UpdateDrink(t *testing.T) {
	db, err := pgxpool.New(context.TODO(), "host=localhost user=username password=password dbname=dbname sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}
	_, err = db.Exec(context.TODO(), QUERY_CREATE_TABLES)
	defer db.Exec(context.TODO(), QUERY_DROP_TABLES)
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

	err = db.QueryRow(ctx, "SELECT name, tags FROM drinks WHERE name = $1", drink.Name).Scan(&updatedDrink.name, &updatedDrink.tags)
	if err != nil {
		t.Fatalf("Failed to retrieve updated drink: %v", err)
	}
	resDrink := fromModelToController(updatedDrink)
	if !reflect.DeepEqual(drink, resDrink) {
		t.Errorf("Updated drink does not match: got %v, want %v", resDrink, drink)
	}
}

func TestSQLDrinkModel_DeleteDrink(t *testing.T) {
	db, err := pgxpool.New(context.TODO(), "host=localhost user=username password=password dbname=dbname sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()
	if err != nil {
		t.Fatalf("Failed to drop table: %v", err)
	}
	_, err = db.Exec(context.TODO(), QUERY_CREATE_TABLES)
	defer db.Exec(context.TODO(), QUERY_DROP_TABLES)
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

	err = db.QueryRow(ctx, "SELECT name, tags FROM drinks WHERE name = $1", drink.Name).Scan(&deletedDrink.name, &deletedDrink.tags)
	if err == nil {
		t.Fatalf("Deleted drink was found: %v", err)
	}
	err = model.DeleteDrink(ctx, drink.Name)
	if err != ErrNotFound {
		t.Fatalf("Error introspection failed: %v", err)
	}
}

func TestSQLDrinkModel_DrinksByTags(t *testing.T) {
	db, err := pgxpool.New(context.TODO(), "host=localhost user=username password=password dbname=dbname sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()
	_, err = db.Exec(context.TODO(), QUERY_CREATE_TABLES)
	defer db.Exec(context.TODO(), QUERY_DROP_TABLES)
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
	db, err := pgxpool.New(context.TODO(), "host=localhost user=username password=password dbname=dbname sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}
	defer db.Close()
	_, err = db.Exec(context.TODO(), QUERY_CREATE_TABLES)
	defer db.Exec(context.TODO(), QUERY_DROP_TABLES)
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
