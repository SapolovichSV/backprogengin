package model

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

func (t tags) Value() (driver.Value, error) {
	var tags []string
	for _, v := range t {
		tags = append(tags, v.Name)
	}
	return strings.Join(tags, ","), nil
}
func (t *tags) Scan(src interface{}) error {
	stringTags, err := src.(string)
	if err {
		return fmt.Errorf("could not convert %v to string", src)
	}
	*t = tags{}
	strings := strings.Split(stringTags, "")
	for _, v := range strings {
		*t = append(*t, tag{Name: v})
	}
	return nil
}
