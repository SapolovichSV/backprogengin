package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type tags []tag

type tag struct {
	Name string
}

// Value implements the driver.Valuer interface for tags
func (t tags) Value() (driver.Value, error) {
	var tagNames []string
	for _, v := range t {
		tagNames = append(tagNames, v.Name)
	}
	return strings.Join(tagNames, ","), nil
}

// Scan implements the sql.Scanner interface for tags
func (t *tags) Scan(src interface{}) error {
	stringTags, ok := src.(string)
	if !ok {
		return fmt.Errorf("could not convert %v to string", src)
	}
	*t = tags{}
	for _, v := range strings.Split(stringTags, ",") {
		*t = append(*t, tag{Name: v})
	}
	return nil
}
