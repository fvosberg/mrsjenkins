package todo

import (
	"encoding/json"
	"net/http"
)

type Todo struct {
	Title       string `validate:"nonzero",json:"title"`
	Description string `validate:"-",json:"description"`
}

func newFromRequest(r *http.Request) (Todo, error) {
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	return todo, err
}
