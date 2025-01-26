package model

import (
	"context"

	"github.com/SapolovichSV/backprogeng/internal/user/entities"
	"github.com/SapolovichSV/backprogeng/internal/user/model/queries"
	"github.com/SapolovichSV/backprogeng/internal/user/model/validate"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SQLUserModel struct {
	db *pgxpool.Pool
}

//go:generate mockgen -destination=../mocks/user/mock_user.go -package=mocks github.com/SapolovichSV/backprogeng/internal/user/model userModel
type userModel interface {
	CreateUser(ctx context.Context, user entities.User) (entities.User, error)
	UserByID(ctx context.Context, id int) (entities.User, error)
	AddFav(ctx context.Context, drinkName string, useriD int) (entities.User, error)
}

func New(db *pgxpool.Pool) *SQLUserModel {
	return &SQLUserModel{
		db: db,
	}
}

func (m *SQLUserModel) CreateUser(ctx context.Context, user entities.User) (entities.User, error) {
	if err := validate.UserName(user.Username); err != nil {
		return entities.User{}, err
	}
	if err := validate.VPassword(user.Password); err != nil {
		return entities.User{}, err
	}
	query := queries.New(m.db, ctx)
	drinksId, err := query.DrinksIdByDrinkNames(user.FavouritesDrinkName)
	if err != nil {
		return entities.User{}, err
	}
	user.ID, err = query.CreateUser(user.Username, user.Password)
	if err != nil {
		return entities.User{}, err
	}
	for _, drinkID := range drinksId {
		query.AddToUserNewFavoriteDrink(user.ID, drinkID)
	}
	return user, nil
}

func (m *SQLUserModel) UserByID(ctx context.Context, id int) (entities.User, error) {
	userRes := entities.User{}
	query := queries.New(m.db, ctx)
	userRes, err := query.UserWithHisFavsByUserID(id)
	if err != nil {
		return entities.User{}, err
	}
	userRes.ID = id
	return userRes, nil
}
func (m *SQLUserModel) AddFav(ctx context.Context, drinkName string, userID int) (res entities.User, err error) {
	query := queries.New(m.db, ctx)
	drinkID, err := query.DrinkIDByName(drinkName)
	if err != nil {
		return entities.User{}, err
	} // err == nil
	res, err = query.UserWithHisFavsByUserID(userID)
	if err != nil {
		return entities.User{}, err
	}
	if err := query.AddToUserNewFavoriteDrink(userID, drinkID); err != nil {
		return entities.User{}, err
	}
	res.FavouritesDrinkName = append(res.FavouritesDrinkName, drinkName)
	res.ID = userID
	return res, nil
}
