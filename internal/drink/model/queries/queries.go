package queries

import (
	"context"

	"github.com/SapolovichSV/backprogeng/internal/drink/entities"
	"github.com/SapolovichSV/backprogeng/internal/errlib"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Query struct {
	ctx context.Context
	db  *pgxpool.Pool
}

const TABLE_NAME = "drinks"

func New(ctx context.Context, db *pgxpool.Pool) *Query {
	return &Query{
		ctx: ctx,
		db:  db,
	}
}
func (q *Query) DrinkByName(name string) (entities.Drink, error) {
	sql := `SELECT id,name,tags FROM drinks
	WHERE name = $1`
	var drink entities.Drink
	err := q.db.QueryRow(q.ctx, sql, name).Scan(&drink.ID, &drink.Name, &drink.Tags)
	if err != nil {
		return entities.Drink{}, errlib.WrapError(err, TABLE_NAME, "drink")
	}
	return drink, nil
}
func (q *Query) CreateDrink(drinkName string) (entities.Drink, error) {

	sql := `INSERT INTO drinks
	(name)
	VALUES($1);`
	_, err := q.db.Exec(q.ctx, sql, drinkName)
	if err != nil {
		return entities.Drink{}, errlib.WrapError(err, "drinks", "drink can't be created")
	}
	var resultDrink entities.Drink
	sql = `SELECT id,name
	FROM drinks
	WHERE name=$1;`
	err = q.db.QueryRow(q.ctx, sql, drinkName).Scan(&resultDrink.ID, &resultDrink.Name)
	if err != nil {
		return entities.Drink{}, errlib.WrapError(err, "drinks", "drink was created but not found")
	}
	return resultDrink, nil
}
func (q *Query) SetTagsToDrink(drinkname string, tags tags) (entities.Drink, error) {

	sql := `UPDATE drinks
	SET tags = $1
	WHERE name = $2;`
	_, err := q.db.Exec(q.ctx, sql, tags, drinkname)
	if err != nil {
		return entities.Drink{}, errlib.WrapError(err, "drinks", "tags can't be set to drink")
	}

	var resultDrink entities.Drink
	haveTags := ToTags([]string{})

	sql = `SELECT id,name,tags
	FROM drinks
	WHERE name=$1;`
	err = q.db.QueryRow(q.ctx, sql, drinkname).Scan(&resultDrink.ID, &resultDrink.Name, &haveTags)
	resultDrink.Tags = FromTags(haveTags)
	if err != nil {
		return entities.Drink{}, errlib.WrapError(err, "drinks", "tags was set but drink not found")
	}

	return resultDrink, nil
}
