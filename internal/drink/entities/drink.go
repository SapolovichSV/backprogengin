package entities

type Drink struct {
	ID   int      `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}
