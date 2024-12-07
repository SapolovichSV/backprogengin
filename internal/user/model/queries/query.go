package queries

import (
	"context"

	"github.com/SapolovichSV/backprogeng/internal/user/entities"
	"github.com/SapolovichSV/backprogeng/internal/user/model/modelerrors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Query struct {
	db  *pgxpool.Pool
	ctx context.Context
}

func New(db *pgxpool.Pool, ctx context.Context) *Query {
	return &Query{
		db:  db,
		ctx: ctx,
	}
}
func (q *Query) DrinksIdByDrinkNames(drinknames entities.Drinknames) ([]int, error) {
	drinksId := make([]int, len(drinknames))
	for i, drinkName := range drinknames {
		sql := `SELECT id 
		FROM drinks
		WHERE name=$1;
		`
		err := q.db.QueryRow(q.ctx, sql, drinkName).Scan(&drinksId[i])
		if err == pgx.ErrNoRows {
			return nil, modelerrors.NotFoundErr{Where: "drinks", What: drinkName}
		} else if err != nil {
			return nil, modelerrors.UnexpectedErr{Where: "drinks", Err: err}
		}
	}
	return drinksId, nil
}
func (q *Query) UserWithHisFavsByUserID(id int) (entities.User, error) {
	queryUser := `SELECT users.username,users.password,drinks.name
        FROM users INNER JOIN drinks ON drinks.id IN (SELECT favs.drink_id
	FROM favs
	WHERE user_id = $1)
WHERE users.id = $2;`
	var res entities.User
	row, _ := q.db.Query(q.ctx, queryUser, id, id)

	for i := 0; row.Next(); i++ {
		var drinkName string
		if i == 0 {
			row.Scan(&res.Username, &res.Password, &drinkName)
		} else {
			row.Scan(nil, nil, &drinkName)
		}
		res.FavouritesDrinkName = append(res.FavouritesDrinkName, drinkName)
	}
	if row.Err() != nil {
		return entities.User{}, modelerrors.WrapError(row.Err(), "users", "user")
	} else if len(res.Username) == 0 {
		return entities.User{}, modelerrors.NotFoundErr{Where: "users", What: "user"}
	}
	return res, nil
}
func (q *Query) AddToUserNewFavoriteDrink(userID int, drinkID int) error {
	queryAddToUserNewFavDrink := `INSERT INTO favs (user_id,drink_id)
	VALUES ($1,$2);`
	if _, err := q.db.Exec(q.ctx, queryAddToUserNewFavDrink, userID, drinkID); err != nil {
		return modelerrors.WrapError(err, "favs", "cannot add new favorite drink")
	}
	return nil
}
func (q *Query) DrinkIDByName(drinkname string) (int, error) {
	queryGetDrinkID := `SELECT drinks.id
	FROM drinks
	WHERE drinks.name = $1;
	`
	var drinkID int
	if err := q.db.QueryRow(q.ctx, queryGetDrinkID, drinkname).Scan(&drinkID); err != nil {
		return 0, modelerrors.WrapError(err, "drinks", drinkname)
	}
	return drinkID, nil
}
func (q *Query) CreateUser(username string, password string) (int, error) {
	sql := "INSERT INTO users (username,password) VALUES ($1,$2) RETURNING id"
	var userID int
	err := q.db.QueryRow(q.ctx, sql, username, password).Scan(&userID)
	if err != nil {
		return 0, modelerrors.WrapError(err, "users", "cannot create user")
	}
	return userID, nil
}
