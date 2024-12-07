package model

import (
	"context"
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

type TestInit struct{}

func NewTest() *TestInit {
	return &TestInit{}
}
func (t *TestInit) Init() (context.Context, *pgxpool.Pool, error) {
	ctx := context.Background()
	db, err := pgxpool.New(context.TODO(), "host=localhost user=username password=password dbname=dbname sslmode=disable")
	if err != nil {
		return nil, nil, err
	}
	_, err = db.Exec(ctx, QUERY_CREATE_TABLES)
	if err != nil {
		return nil, nil, err
	}
	return ctx, db, nil
}
func TestSQLUserModel_CreateUser(t *testing.T) {
	ctx, db, err := NewTest().Init()
	assert.NoError(t, err)
	defer db.Close()
	defer db.Exec(ctx, QUERY_DROP_TABLES)
	tests := []struct {
		name                       string
		beforeCreate               entities.User
		wantErr                    bool
		mustCreateDrinksBeforeTest bool
		drinks                     []drEnt.Drink
	}{
		{
			name: "NoSuchDrinks",
			beforeCreate: entities.User{
				Username:            "Stas002",
				Password:            "amahasla",
				FavouritesDrinkName: entities.Drinknames{"cola", "pussy"},
			},
			wantErr:                    true,
			mustCreateDrinksBeforeTest: false,
			drinks:                     nil,
		},
		{
			name: "SimpleTest",
			beforeCreate: entities.User{
				Username:            "Stas001",
				Password:            "amahasla",
				FavouritesDrinkName: entities.Drinknames{"cola", "pussy"},
			},
			wantErr:                    false,
			mustCreateDrinksBeforeTest: true,
			drinks: []drEnt.Drink{
				{Name: "cola", Tags: []string{"spicy", "spice"}},
				{Name: "pussy", Tags: []string{"sweet", "spyce"}},
			},
		},
		{
			name: "BadUserName",
			beforeCreate: entities.User{
				Username:            "",
				Password:            "amahasla",
				FavouritesDrinkName: nil,
			},
			wantErr:                    true,
			mustCreateDrinksBeforeTest: false,
			drinks:                     nil,
		},
		{
			name: "BadPassword",
			beforeCreate: entities.User{
				Username:            "Stas003",
				Password:            "",
				FavouritesDrinkName: nil,
			},
			wantErr:                    true,
			mustCreateDrinksBeforeTest: false,
		},
	}
	drinkModel := model.New(db)
	m := New(db)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mustCreateDrinksBeforeTest {
				for _, drink := range tt.drinks {
					drinkModel.CreateDrink(ctx, drink)
				}
			}
			have, err := m.CreateUser(ctx, tt.beforeCreate)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				tt.beforeCreate.ID = have.ID
				assert.Equal(t, tt.beforeCreate, have)
			}

		})
	}
}

func TestSQLUserModel_UserByID(t *testing.T) {
	ctx, db, err := NewTest().Init()
	assert.NoError(t, err)
	defer db.Close()
	defer db.Exec(ctx, QUERY_DROP_TABLES)
	tests := []struct {
		name                       string
		want                       entities.User
		wantErr                    bool
		mustCreateUserBeforeTest   bool
		mustCreateDrinksBeforeTest bool
		drinks                     []drEnt.Drink
	}{
		{
			name: "SimpleTest",
			want: entities.User{
				ID:                  0,
				Username:            "Stas228",
				Password:            "amahasla",
				FavouritesDrinkName: []string{"cola", "pussy"},
			},
			wantErr:                    false,
			mustCreateUserBeforeTest:   true,
			mustCreateDrinksBeforeTest: true,
			drinks: []drEnt.Drink{
				{Name: "cola", Tags: []string{"spicy", "spice"}},
				{Name: "pussy", Tags: []string{"sweet", "spyce"}},
			},
		},
		{
			name: "UserNotExist",
			want: entities.User{
				ID: 228,
			},
			wantErr:                    true,
			mustCreateUserBeforeTest:   false,
			mustCreateDrinksBeforeTest: false,
			drinks:                     nil,
		},
	}
	drinkModel := model.New(db)
	m := New(db)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mustCreateUserBeforeTest {
				if tt.mustCreateDrinksBeforeTest {
					for _, drink := range tt.drinks {
						drinkModel.CreateDrink(ctx, drink)
					}
				}
				user, err := m.CreateUser(ctx, tt.want)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					tt.want.ID = user.ID
					assert.Equal(t, tt.want, user)
				}
			} else {
				user, err := m.UserByID(ctx, tt.want.ID)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.want, user)
				}
			}
		})
	}
}

func TestSQLUserModel_AddFav(t *testing.T) {
	ctx, db, err := NewTest().Init()
	assert.NoError(t, err)
	defer db.Close()
	defer db.Exec(ctx, QUERY_DROP_TABLES)
	tests := []struct {
		name                       string
		beforeAdd                  entities.User
		want                       entities.User
		wantErr                    bool
		mustCreateUserBeforeTest   bool
		mustCreateDrinksBeforeTest bool
		drinks                     []drEnt.Drink
		addingDrinkName            string
	}{
		{
			name: "SimpleTest",
			beforeAdd: entities.User{
				Username:            "Stas228",
				Password:            "amahasla",
				FavouritesDrinkName: []string{"cola", "pussy"},
			},
			want: entities.User{
				Username:            "Stas228",
				Password:            "amahasla",
				FavouritesDrinkName: []string{"cola", "pussy", "pepsi"},
			},
			wantErr:                    false,
			mustCreateUserBeforeTest:   true,
			mustCreateDrinksBeforeTest: true,
			drinks: []drEnt.Drink{
				{Name: "cola", Tags: []string{"spicy", "spice"}},
				{Name: "pussy", Tags: []string{"sweet", "spyce"}},
				{Name: "pepsi", Tags: []string{"sweet", "spyce"}},
			},
			addingDrinkName: "pepsi",
		},
	}
	drinkModel := model.New(db)
	m := New(db)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mustCreateUserBeforeTest {
				if tt.mustCreateDrinksBeforeTest {
					for _, drink := range tt.drinks {
						drinkModel.CreateDrink(ctx, drink)
					}
				}
				createdUser, err := m.CreateUser(ctx, tt.beforeAdd)
				assert.NoError(t, err)

				have, err := m.AddFav(ctx, tt.addingDrinkName, createdUser.ID)
				tt.want.ID = createdUser.ID
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.want, have)
				}
			} else {
				if tt.mustCreateDrinksBeforeTest {
					for _, drink := range tt.drinks {
						drinkModel.CreateDrink(ctx, drink)
					}
				}
				_, err := m.AddFav(ctx, tt.addingDrinkName, tt.beforeAdd.ID)
				if tt.wantErr {
					assert.Error(t, err)
				}
			}
		})
	}
}
