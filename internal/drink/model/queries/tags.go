package queries

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

type tags []tag

type tag struct {
	Name string
}

func (t tags) Value() (driver.Value, error) {
	var tagNames []string
	for _, v := range t {
		tagNames = append(tagNames, v.Name)
	}
	return strings.Join(tagNames, ","), nil
}

func (t *tags) Scan(src interface{}) error {
	switch v := src.(type) {
	case string:
		// Разбиваем строку на теги
		*t = tags{}
		for _, name := range strings.Split(v, ",") {
			*t = append(*t, tag{Name: strings.TrimSpace(name)})
		}
		return nil
	case []byte:
		// Если значение приходит в виде []byte
		return t.Scan(string(v))
	default:
		return fmt.Errorf("could not convert %T to tags", src)
	}
}
func ToTags(tgs []string) tags {
	var t tags
	for _, v := range tgs {
		t = append(t, tag{Name: v})
	}
	return t
}
func FromTags(t tags) []string {
	var tgs []string
	for _, v := range t {
		tgs = append(tgs, v.Name)
	}
	return tgs
}
