package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict = errors.New("edit conflict")
)

type Models struct {
	Reviews interface {
		Insert(review *Review) error
		Get(id int64) (*Review, error)
		Update(review *Review) error
		Delete(id int64) error
	}
	Pizzas interface {
		Insert(pizza *Pizza) error
		Get(id int64) (*Pizza, error)
		Update(pizza *Pizza) error
		Delete(id int64) error
		GetAll() ([]*Pizza, error)
	}
	Images interface {
		Insert(image *Image) error
		Get(id int64) (*Image, error)
		Update(image *Image) error
		Delete(id int64) error
	}
	Venues interface {
		Insert(venue *Venue) error
		Get(id int64) (*Venue, error)
		Update(venue *Venue) error
		Delete(id int64) error
		GetAll() ([]*Venue, error)
	}
	VenuePizzas interface {
		Insert(venuePizza *VenuePizza) error
		Get(id int64) (*VenuePizza, error)
		Update(venuePizza *VenuePizza) error
		Delete(id int64) error
		GetAll() ([]*VenuePizzaMixin, error)
	}

}

func NewModels(db *sql.DB) Models {
	return Models{
		Reviews: ReviewModel{DB: db},
		Pizzas: PizzaModel{DB: db},
		Images: ImageModel{DB: db},
		Venues: VenueModel{DB: db},
		VenuePizzas: VenuePizzaModel{DB: db},
	}
}

func NewMockModels() Models {
	return Models{
		Reviews: MockReviewModel{},
		Pizzas: MockPizzaModel{},
		Images: MockImageModel{},
		Venues: MockVenueModel{},
		VenuePizzas: MockVenuePizzaModel{},
	}
}