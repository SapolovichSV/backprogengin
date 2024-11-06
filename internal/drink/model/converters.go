package model

import (
	"fmt"

	"github.com/SapolovichSV/backprogeng/internal/drink/entities"
)

func fromControllerToModel(c entities.Drink) Drink {
	return Drink{
		name: c.Name,
		tags: fromControllerToModelTags(c.Tags),
	}
}
func fromModelToController(m Drink) entities.Drink {
	return entities.Drink{
		Name: m.name,
		Tags: fromModelToControllerTags(m.tags),
	}
}
func fromModelToControllerTags(m tags) []string {
	var t []string
	for _, v := range m {
		t = append(t, v.Name)
	}
	return t
}
func fromControllerToModelTags(c []string) tags {
	var t tags
	for _, v := range c {
		t = append(t, tag{Name: v})
	}
	return t
}
func wrapifErrorInModel(msg string, err error) error {
	if err != nil {
		return fmt.Errorf("%s : %s", msg, err.Error())
	}
	return nil
}
