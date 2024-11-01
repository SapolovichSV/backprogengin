package model

import "github.com/SapolovichSV/backprogeng/internal/drink/controller"

func fromControllerToModel(c controller.Drink) Drink {
	return Drink{
		name: c.Name,
		tags: fromControllerToModelTags(c.Tags),
	}
}
func fromModelToController(m Drink) controller.Drink {
	return controller.Drink{
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
