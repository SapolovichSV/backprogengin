package model

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SapolovichSV/backprogeng/internal/errlib"
	"github.com/SapolovichSV/backprogeng/internal/user/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB:
// ID username password favourites
var ErrNotFound = errors.New("not found")

//Переделить таблицу в нормализованный вид
// ID username password
//another table:
//UserID Fav
//Example
//0 tasty
//0 strong
//1 tasty

type SQLUserModel struct {
	db *pgxpool.Pool
}
type userModel interface {
	CreateUser(ctx context.Context, user entities.User) (entities.User, error)
	UserById(ctx context.Context, id int) (entities.User, error)
	AddFav(ctx context.Context, drinkName string, user entities.User)
}

func New(db *pgxpool.Pool) *SQLUserModel {
	return &SQLUserModel{
		db: db,
	}
}

func (m *SQLUserModel) CreateUser(ctx context.Context, user entities.User) (entities.User, error) {
	sql := "INSERT INTO users (username,password) VALUES ($1,$2) RETURNING id"

	err := m.db.QueryRow(ctx, sql, user.Username, user.Password).Scan(&user.ID)
	if err != nil {
		return entities.User{}, err
	}
	drinksId := make([]int, len(user.FavouritesDrinkName))
	for i, drinkName := range user.FavouritesDrinkName {
		sql := `SELECT id 
		FROM drinks
		WHERE name=$1;
		`
		err := m.db.QueryRow(ctx, sql, drinkName).Scan(&drinksId[i])
		if err == pgx.ErrNoRows {
			return entities.User{}, ErrNotFound
		} else if err != nil {
			return entities.User{}, err
		}
	}
	for _, drinkID := range drinksId {
		sql := "INSERT INTO favs (user_id,drink_id) VALUES ($1,$2)"
		_, err := m.db.Exec(ctx, sql, user.ID, drinkID)
		if err != nil {
			return entities.User{}, err
		}
	}
	return user, nil
}
func (m *SQLUserModel) UserByID(ctx context.Context, id int) (entities.User, error) {
	userRes := entities.User{}
	queryUser := `SELECT users.username,users.password,drinks.name
        FROM users INNER JOIN drinks ON drinks.id IN (SELECT favs.drink_id
	FROM favs
	WHERE user_id = $1)
WHERE users.id = $2;
`
	row, _ := m.db.Query(ctx, queryUser, id, id)

	for i := 0; row.Next(); i++ {
		var drinkName string
		if i == 0 {
			row.Scan(&userRes.Username, &userRes.Password, &drinkName)
		} else {
			row.Scan(nil, nil, &drinkName)
		}
		userRes.FavouritesDrinkName = append(userRes.FavouritesDrinkName, drinkName)
	}
	if userRes.Username == "" || row.Err() == pgx.ErrNoRows {
		return userRes, ErrNotFound
	} else if row.Err() != nil {
		return entities.User{}, nil
	}
	userRes.ID = id
	return userRes, row.Err()
}
func (m *SQLUserModel) AddFav(ctx context.Context, drinkName string, user entities.User) (res entities.User, err error) {
	query := "SELECT id,username,password,favourites FROM users WHERE id = $1"

	err = m.db.QueryRow(ctx, query, user.ID).Scan(&user.ID, &user.Username, &user.Password, &user.FavouritesDrinkName)
	if err == sql.ErrNoRows {
		return entities.User{}, ErrNotFound
	} else if err != nil {
		return entities.User{}, errlib.WrapErr(err, "error getting user")
	}

	user.FavouritesDrinkName = append(user.FavouritesDrinkName, drinkName)
	query = "UPDATE users SET favourites = $1 WHERE id = $2"
	_, err = m.db.Exec(ctx, query, user.FavouritesDrinkName, user.ID)
	if err == sql.ErrNoRows {
		return entities.User{}, ErrNotFound
	} else if err != nil {
		err = errlib.WrapErr(err, "error adding favourite")
	}
	return user, err
}
