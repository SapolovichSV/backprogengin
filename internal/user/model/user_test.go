package model

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	drEnt "github.com/SapolovichSV/backprogeng/internal/drink/entities"
	"github.com/SapolovichSV/backprogeng/internal/drink/model"
	"github.com/SapolovichSV/backprogeng/internal/user/entities"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
)

// sudo docker run --rm --name test-postgres -e POSTGRES_PASSWORD=password -e POSTGRES_USER=username -e POSTGRES_DB=dbname -p 5432:5432 -d postgres
const QUERY_CREATE_TABLES = `CREATE TABLE drinks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    tags TEXT
);
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    password VARCHAR(255)
);
CREATE TABLE favs (
    user_id INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)  ON DELETE CASCADE,
    drink_id INT NOT NULL,
   FOREIGN KEY (drink_id) REFERENCES drinks(id) ON DELETE CASCADE
);`
const QUERY_DROP_TABLES = `DROP TABLE drinks CASCADE;
DROP TABLE users CASCADE;
DROP TABLE favs CASCADE;`

func TestSQLUserModel_CreateUser(t *testing.T) {
	ctx := context.Background()
	db, err := pgxpool.New(context.TODO(), "host=localhost user=username password=password dbname=dbname sslmode=disable")
	assert.NoError(t, err, "error db connect")
	_, err = db.Exec(ctx, QUERY_CREATE_TABLES)
	assert.NoError(t, err, "error db connect")
	defer db.Close()
	defer db.Exec(ctx, QUERY_DROP_TABLES)
	tests := []struct {
		name         string
		beforeCreate entities.User
		wantErr      bool
	}{
		{
			name: "test01",
			beforeCreate: entities.User{
				Username:            "Stas228",
				Password:            "amahasla",
				FavouritesDrinkName: entities.Drinknames{"cola", "pussy"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(db)
			afterCreate, err := m.CreateUser(ctx, tt.beforeCreate)
			if !tt.wantErr {
				assert.NoError(t, err)
				assert.Equal(t, tt.beforeCreate, afterCreate)
			} else {
				assert.EqualError(t, err, ErrNotFound.Error(), "wrong err")
			}
		})
	}
	tests[0].wantErr = false

	drinkModel := model.New(db)
	drinkModel.CreateDrink(ctx, drEnt.Drink{Name: "cola", Tags: []string{"spicy", "spice"}})
	drinkModel.CreateDrink(ctx, drEnt.Drink{Name: "pussy", Tags: []string{"sweet", "spyce"}})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New(db)
			afterCreate, err := m.CreateUser(ctx, tt.beforeCreate)
			if !tt.wantErr {
				assert.NoError(t, err)
				tt.beforeCreate.ID = afterCreate.ID
				assert.Equal(t, tt.beforeCreate, afterCreate)
			} else {
				assert.EqualError(t, err, ErrNotFound.Error(), "wrong err")
			}
		})
	}
	fmt.Print()
}

func TestSQLUserModel_UserByID(t *testing.T) {
	ctx := context.Background()
	db, err := pgxpool.New(context.TODO(), "host=localhost user=username password=password dbname=dbname sslmode=disable")
	assert.NoError(t, err, "error db connect")
	_, err = db.Exec(ctx, QUERY_CREATE_TABLES)
	assert.NoError(t, err, "error db connect")
	defer db.Close()
	defer db.Exec(ctx, QUERY_DROP_TABLES)
	tests := []struct {
		name    string
		want    entities.User
		wantErr bool
	}{
		{
			name: "Test01",
			want: entities.User{
				ID:                  0,
				Username:            "Stas228",
				Password:            "amahasla",
				FavouritesDrinkName: []string{"cola", "pussy"},
			},
		},
	}
	testCase1 := tests[0]
	user1 := testCase1.want
	m := New(db)
	drinkModel := model.New(db)
	drinkModel.CreateDrink(ctx, drEnt.Drink{Name: "cola", Tags: []string{"spicy", "spice"}})
	drinkModel.CreateDrink(ctx, drEnt.Drink{Name: "pussy", Tags: []string{"sweet", "spyce"}})
	want, err := m.CreateUser(ctx, user1)
	assert.NoError(t, err)
	get, err := m.UserByID(ctx, want.ID)
	assert.NoError(t, err)
	assert.Equal(t, want, get)
}

func TestSQLUserModel_AddFav(t *testing.T) {
	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		ctx       context.Context
		drinkName string
		user      entities.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantRes entities.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &SQLUserModel{
				db: tt.fields.db,
			}
			gotRes, err := m.AddFav(tt.args.ctx, tt.args.drinkName, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLUserModel.AddFav() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("SQLUserModel.AddFav() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
