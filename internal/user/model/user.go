package model

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SapolovichSV/backprogeng/internal/errlib"
	"github.com/SapolovichSV/backprogeng/internal/user/entities"
)

// DB:
// ID username password favourites
var ErrNotFound = errors.New("not found")

type SQLUserModel struct {
	db *sql.DB
}

func New(db *sql.DB) *SQLUserModel {
	return &SQLUserModel{
		db: db,
	}
}

func (m *SQLUserModel) CreateUser(ctx context.Context, user entities.User) (entities.User, error) {
	sql := "INSERT INTO users (username,password,favourites) VALUES ($1,$2,$3) RETURNING id"
	err := m.db.QueryRowContext(ctx, sql, user.Username, user.Password, user.FavouritesDrinkName).Scan(&user.ID)
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}
func (m *SQLUserModel) UserByID(ctx context.Context, id int) (entities.User, error) {
	var user entities.User
	query := "SELECT id,username,password,favourites FROM users WHERE id=$1"
	err := m.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Password, &user.FavouritesDrinkName)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.User{}, ErrNotFound
		} else {
			return entities.User{}, err
		}
	}
	return user, nil
}
func (m *SQLUserModel) AddFav(ctx context.Context, drinkName string, user entities.User) (res entities.User, err error) {
	query := "SELECT id,username,password,favourites FROM users WHERE id = $1"

	err = m.db.QueryRowContext(ctx, query, user.ID).Scan(&user.ID, &user.Username, &user.Password, &user.FavouritesDrinkName)
	if err == sql.ErrNoRows {
		return entities.User{}, ErrNotFound
	} else if err != nil {
		return entities.User{}, errlib.WrapErr(err, "error getting user")
	}

	user.FavouritesDrinkName = append(user.FavouritesDrinkName, drinkName)
	query = "UPDATE users SET favourites = $1 WHERE id = $2"
	_, err = m.db.ExecContext(ctx, query, user.FavouritesDrinkName, user.ID)
	if err == sql.ErrNoRows {
		return entities.User{}, ErrNotFound
	} else if err != nil {
		err = errlib.WrapErr(err, "error adding favourite")
	}
	return user, err
}
