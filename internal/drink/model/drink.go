package model

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"strings"

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

func (t tags) Value() (driver.Value, error) {
	var tags []string
	for _, v := range t {
		tags = append(tags, v.Name)
	}
	return strings.Join(tags, ","), nil
}
func (t *tags) Scan(src interface{}) error {
	stringTags, err := src.(string)
	if err {
		return fmt.Errorf("could not convert %v to string", src)
	}
	*t = tags{}
	strings := strings.Split(stringTags, "")
	for _, v := range strings {
		*t = append(*t, tag{Name: v})
	}
	return nil
}

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
func (m *SQLDrinkModel) CreateDrink(dCont controller.Drink) (controller.Drink, error) {
	d := fromControllerToModel(dCont)
	sql, args, err := sq.Insert("drinks").Columns("name", "tags").Values(d.name, d.tags).ToSql()
	if err != nil {
		return fromModelToController(Drink{}), err
	}
	_, err = m.db.Exec(sql, args...)
	if err != nil {
		return fromModelToController(Drink{}), err
	}
	return fromModelToController(d), nil
}
func (m *SQLDrinkModel) UpdateDrink(dCont controller.Drink) (controller.Drink, error) {
	d := fromControllerToModel(dCont)
	sql, args, err := sq.Update("drinks").Set("tags", d.tags).Where(sq.Eq{"name": d.name}).ToSql()
	if err != nil {
		return fromModelToController(Drink{}), err
	}
	_, err = m.db.Exec(sql, args...)
	if err != nil {
		return fromModelToController(Drink{}), err
	}
	return fromModelToController(d), nil
}
func (m *SQLDrinkModel) DeleteDrink(name string) error {
	sql, args, err := sq.Delete("drinks").Where(sq.Eq{"name": name}).ToSql()
	if err != nil {
		return err
	}
	_, err = m.db.Exec(sql, args...)
	if err != nil {
		return err
	}
	return nil
}
func (m *SQLDrinkModel) DrinksByTags(tagsCont []string) ([]controller.Drink, error) {
	tags := fromControllerToModelTags(tagsCont)
	sql, args, err := sq.Select("name", "tags").From("drinks").Where(sq.Eq{"tags": tags}).ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := m.db.Query(sql, args...)
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
func (m *SQLDrinkModel) AllDrinks() ([]controller.Drink, error) {
	sql, args, err := sq.Select("name", "tags").From("drinks").ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := m.db.Query(sql, args...)
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
func (m *SQLDrinkModel) DrinkByName(name string) (controller.Drink, error) {
	sql, args, err := sq.Select("name", "tags").From("drinks").Where(sq.Eq{"name": name}).ToSql()
	if err != nil {
		return controller.Drink{}, err
	}
	row := m.db.QueryRow(sql, args...)
	var d Drink
	err = row.Scan(&d.name, &d.tags)
	if err != nil {
		return controller.Drink{}, err
	}
	return fromModelToController(d), nil
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
