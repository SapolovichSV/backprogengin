package entities

type Drink struct {
	ID   int      `json:"id,omitempty" example:"12"`
	Name string   `json:"name" example:"Coca Cola"`
	Tags []string `json:"tags" example:"[\"soda\",\"cola\"]"`
}
