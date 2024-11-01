package model

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/SapolovichSV/backprogeng/internal/drink/controller"
)

// DB:
// drinks
// id | name | tags
var ErrNotFound = fmt.Errorf("not found")

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
func (m *SQLDrinkModel) CreateDrink(ctx context.Context, dCont controller.Drink) (controller.Drink, error) {
	d := fromControllerToModel(dCont)
	sql, args, err := sq.Insert("drinks").Columns("name", "tags").Values(d.name, d.tags).ToSql()
	if err != nil {
		return fromModelToController(Drink{}), err
	}
	_, err = m.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fromModelToController(Drink{}), err
	}
	return fromModelToController(d), nil
}
func (m *SQLDrinkModel) UpdateDrink(ctx context.Context, dCont controller.Drink) (controller.Drink, error) {
	d := fromControllerToModel(dCont)
	sql, args, err := sq.Update("drinks").Set("tags", d.tags).Where(sq.Eq{"name": d.name}).ToSql()
	if err != nil {
		return fromModelToController(Drink{}), err
	}
	_, err = m.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return fromModelToController(Drink{}), err
	}
	return fromModelToController(d), nil
}
func (m *SQLDrinkModel) DeleteDrink(ctx context.Context, name string) error {
	sql, args, err := sq.Delete("drinks").Where(sq.Eq{"name": name}).ToSql()
	if err != nil {
		return err
	}
	_, err = m.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}
func (m *SQLDrinkModel) DrinksByTags(ctx context.Context, tagsCont []string) ([]controller.Drink, error) {
	tags := fromControllerToModelTags(tagsCont)
	sql, args, err := sq.Select("name", "tags").From("drinks").Where(sq.Eq{"tags": tags}).ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := m.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	var drinks []controller.Drink
	for rows.Next() {
		var d Drink
		err = rows.Scan(&d.name, &d.tags)
		if err != nil {
			return nil, err
		}
		drinks = append(drinks, fromModelToController(d))
	}
	if len(drinks) == 0 {
		return nil, ErrNotFound
	}
	return drinks, nil
}
func (m *SQLDrinkModel) AllDrinks(ctx context.Context, id int) ([]controller.Drink, error) {
	sql, args, err := sq.Select("name", "tags").From("drinks").ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := m.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	var drinks []controller.Drink
	for rows.Next() {
		var d Drink
		err = rows.Scan(&d.name, &d.tags)
		if err != nil {
			return nil, err
		}
		drinks = append(drinks, fromModelToController(d))
	}
	if len(drinks) == 0 {
		return nil, ErrNotFound
	}

	return drinks, nil
}
func (m *SQLDrinkModel) DrinkByName(ctx context.Context, name string) (controller.Drink, error) {
	sql, args, err := sq.Select("name", "tags").From("drinks").Where(sq.Eq{"name": name}).ToSql()
	if err != nil {
		return controller.Drink{}, err
	}
	row := m.db.QueryRowContext(ctx, sql, args...)
	var d Drink
	err = row.Scan(&d.name, &d.tags)
	if err != nil {
		return controller.Drink{}, err
	}
	return fromModelToController(d), nil
}
