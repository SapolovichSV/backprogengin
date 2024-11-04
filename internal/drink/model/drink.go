package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/SapolovichSV/backprogeng/internal/drink/entities"
)

// DB:
// drinks
// id | name | tags
var ErrNotFound = fmt.Errorf("not found")
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

// Drink is a struct that represents a drink
type Drink struct {
	name string
	tags tags
}
type SQLDrinkModel struct {
	db *sql.DB
}
type DrinkModel interface {
	CreateDrink(ctx context.Context, dCont entities.Drink) (entities.Drink, error)
	UpdateDrink(ctx context.Context, dCont entities.Drink) (entities.Drink, error)
	DeleteDrink(ctx context.Context, name string) error
	DrinksByTags(ctx context.Context, tagsCont []string) ([]entities.Drink, error)
	AllDrinks(ctx context.Context, id int) ([]entities.Drink, error)
	DrinkByName(ctx context.Context, name string) (entities.Drink, error)
}

func NewSQLDrinkModel(db *sql.DB) *SQLDrinkModel {
	return &SQLDrinkModel{
		db: db,
	}
}
func (m *SQLDrinkModel) CreateDrink(ctx context.Context, dCont entities.Drink) (entities.Drink, error) {
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
func (m *SQLDrinkModel) UpdateDrink(ctx context.Context, dCont entities.Drink) (entities.Drink, error) {
	d := fromControllerToModel(dCont)
	sql, args, err := sq.Update("drinks").Set("tags", d.tags).Where(squirrel.Eq{"name": d.name}).ToSql()
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
	sql, args, err := sq.Delete("drinks").Where(squirrel.Eq{"name": name}).ToSql()
	if err != nil {
		return err
	}
	_, err = m.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}
	return nil
}
func (m *SQLDrinkModel) DrinksByTags(ctx context.Context, tagsCont []string) ([]entities.Drink, error) {
	tags := fromControllerToModelTags(tagsCont)
	sql, args, err := sq.Select("name", "tags").From("drinks").Where(squirrel.Eq{"tags": tags}).ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := m.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	var drinks []entities.Drink
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
func (m *SQLDrinkModel) AllDrinks(ctx context.Context, id int) ([]entities.Drink, error) {
	sql, args, err := sq.Select("name", "tags").From("drinks").Where(squirrel.Gt{"id": id}).ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := m.db.QueryContext(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	var drinks []entities.Drink
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
func (m *SQLDrinkModel) DrinkByName(ctx context.Context, name string) (entities.Drink, error) {
	sql, args, err := sq.Select("name", "tags").From("drinks").Where(squirrel.Eq{"name": name}).ToSql()
	if err != nil {
		return entities.Drink{}, err
	}
	row := m.db.QueryRowContext(ctx, sql, args...)
	var d Drink
	err = row.Scan(&d.name, &d.tags)
	if err != nil {
		return entities.Drink{}, err
	}
	return fromModelToController(d), nil
}
