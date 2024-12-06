package entities

import (
	"database/sql/driver"
	"errors"
	"strings"
)

type User struct {
	ID                  int        `json:"id"`
	Username            string     `json:"username"`
	Password            string     `json:"password"`
	FavouritesDrinkName Drinknames `json:"drinknames"`
}
type Drinknames []string

func (ds Drinknames) Value() (driver.Value, error) {
	return strings.Join(ds, ","), nil
}
func (ds *Drinknames) Scan(src interface{}) error {
	stringTags, ok := src.(string)
	if !ok {
		return errors.New("could not convert to string")
	}
	*ds = Drinknames{}

	for _, v := range strings.Split(stringTags, ",") {
		*ds = append(*ds, v)
	}
	return nil
}
