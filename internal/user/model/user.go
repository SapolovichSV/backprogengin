package model

import (
	"context"
	"database/sql"
	"errors"

	"github.com/SapolovichSV/backprogeng/internal/user/entities"
)

var ErrNotFound = errors.New("not found")

type SQLUserModel struct {
	db *sql.DB
}

func (m *SQLUserModel) CreateUser(ctx context.Context, user entities.User) (entities.User, error) {
	sql := "INSERT INTO users (username,password) VALUES ($1,$2,$3) RETURNING id"
	err := m.db.QueryRowContext(ctx, sql, user.Username, user.Password).Scan(&user.ID)
	if err != nil {
		return entities.User{}, err
	}
	return user, nil
}
func (m *SQLUserModel) UserByID(ctx context.Context, id int) (entities.User, error) {
	var user entities.User
	query := "SELECT id,username,password FROM users WHERE id=$1"
	err := m.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.User{}, ErrNotFound
		} else {
			return entities.User{}, err
		}
	}
	return user, nil
}
