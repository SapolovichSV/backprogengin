package model

import (
	"context"
	"errors"

	"github.com/SapolovichSV/backprogeng/internal/user/entities"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

const QUERY_IDS_OF_FAV_DRINKS_USERS = `(SELECT favs.drink_id
	FROM favs
	WHERE user_id = $1)`

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
func (m *SQLUserModel) AddFav(ctx context.Context, drinkName string, userID int) (res entities.User, err error) {
	queryUser := `SELECT users.username,users.password,drinks.name
        FROM users INNER JOIN drinks ON drinks.id IN (SELECT favs.drink_id
	FROM favs
	WHERE user_id = $1)
WHERE users.id = $2;
`
	queryGetDrinkID := `SELECT drinks.id
	FROM drinks
	WHERE drinks.name = $1;
	`
	var drinkID int
	if err = m.db.QueryRow(ctx, queryGetDrinkID, drinkName).Scan(&drinkID); err != nil {
		return entities.User{}, ErrNotFound
	}
	//Получаю инфу о юзере(нейм,пасс,фаворитес)
	row, _ := m.db.Query(ctx, queryUser, userID, userID)

	for i := 0; row.Next(); i++ {
		var drinkName string
		if i == 0 {
			row.Scan(&res.Username, &res.Password, &drinkName)
		} else {
			row.Scan(nil, nil, &drinkName)
		}
		res.FavouritesDrinkName = append(res.FavouritesDrinkName, drinkName)
	}
	if res.Username == "" || row.Err() == pgx.ErrNoRows {
		return res, ErrNotFound
	} else if row.Err() != nil {
		return entities.User{}, row.Err()
	}
	queryAddToUserNewFavDrink := `INSERT INTO favs (user_id,drink_id)
	VALUES ($1,$2);`
	if _, err := m.db.Exec(ctx, queryAddToUserNewFavDrink, userID, drinkID); err != nil {
		return entities.User{}, err
	}
	res.FavouritesDrinkName = append(res.FavouritesDrinkName, drinkName)
	res.ID = userID
	return res, nil
}
