package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Pizzas interface {
		Insert(pizza *Pizza) error
		Get(id int64) (*Pizza, error)
		Update(pizza *Pizza) error
		Delete(id int64) error
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Pizzas: PizzaModel{DB: db},
	}
}

func NewMockModels() Models {
	return Models{
		Pizzas: MockPizzaModel{},
	}
}