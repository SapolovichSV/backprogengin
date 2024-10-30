package model

import (
	"database/sql"

	"github.com/SapolovichSV/backprogeng/internal/drink/controller"
)

// Drink is a struct that represents a drink
type Drink struct {
	name string
	tags tags
}
type tags []tag
type tag struct {
	Name string
}
type SQLDrinkModel struct {
	db *sql.DB
}

func NewSQLDrinkModel(db *sql.DB) *SQLDrinkModel {
	return &SQLDrinkModel{
		db: db,
	}
}
func (m *SQLDrinkModel) CreateDrink(d Drink) (Drink, error) {
	return Drink{}, nil
}
func (m *SQLDrinkModel) UpdateDrink(d Drink) (Drink, error) {
	return Drink{}, nil
}
func (m *SQLDrinkModel) DeleteDrink(name string) error {
	return nil
}
func (m *SQLDrinkModel) DrinksByTags(tags []string) ([]Drink, error) {
	return nil, nil
}
func (m *SQLDrinkModel) AllDrinks() ([]Drink, error) {
	return nil, nil
}
func (m *SQLDrinkModel) DrinkByName(name string) (Drink, error) {
	return Drink{}, nil
}
func fromControllerToModel(c controller.Drink) Drink {
	return Drink{
		name: c.Name,
		tags: fromControllerToModelTags(c.Tags),
	}
}
func fromControllerToModelTags(c []string) tags {
	var t tags
	for _, v := range c {
		t = append(t, tag{Name: v})
	}
	return t
}
func fromModelToController(m Drink) controller.Drink {
	return controller.Drink{
		Name: m.name,
		Tags: fromModelToControllerTags(m.tags),
	}
}
func fromModelToControllerTags(m tags) []string {
	var t []string
	for _, v := range m {
		t = append(t, v.Name)
	}
	return t
}
