package data

import (
	"time"
	"database/sql"
	_ "errors"
	"context"
	_ "fmt"

	"github.com/tclohm/project-pizza/internal/validator"

	_ "github.com/lib/pq"
)

type VenuePizza struct {
	ID int64 `json:"id"`
	venueId int64 `json:"venue_id"`
	pizzaId int64 `json:"pizza_id"`
}

type VenuePizzaModel struct {
	DB *sql.DB
}

func (vpm VenuePizzaModel) Insert(venuePizza *VenuePizza) error {
	query := `
	INSERT INTO venuepizzas (
		venue_id,
		pizza_id,
	)
	VALUES ($1, $2)
	RETURNING id
	`
	// args slices containing values for the placeholder parameters from the venue struct
	args := []interface{}{
		venuePizza.venue_id, venuePizza.pizza_id
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	// passing in the slice and scanning the system generated id
	return vpm.DB.QueryRowContext(ctx, query, args...).Scan(&venuePizza.ID)
}

func (vpm VenuePizzaModel) Get(id int64) (*VenuePizza, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, 
		venue_id,
		pizza_id,
	FROM venuepizzas WHERE id = $1
	`

	var venuepizza VenuePizza
	// 3-second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	// release resources associated with context before Get() is returned
	// memory leak avoided
	defer cancel()

	err := vpm.DB.QueryRowContext(ctx, query, id).Scan(
		&venuepizza.ID,
		&venuepizza.VenueId,
		&venuepizza.PizzaId
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &venuepizza, nil
}

func (vpm VenuePizzaModel) Update(venuePizza *VenuePizza) error {
	query := `
	UPDATE venuepizzas
	SET venue_id = $1,
		pizza_id = $2,
	WHERE id = $3
	RETURNING id
	`

	args := []interface{}{
		venuePizza.VenueId,
		venuePizza.PizzaId,
		venuePizza.ID
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()
	// query and scan the new value in
	err := vpm.DB.QueryRowContext(ctx, query, args...).Scan(&venuePizza.ID)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (vpm VenuePizzaModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM venuepizzas
		WHERE id = $1
	`

	result, err := vpm.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrRecordNotFound
	}
	
	return nil
}


type MockVenuePizzaModel struct {}

func (vpm MockVenuePizzaModel) Insert(venue *Venue) error {
	return nil
}

func (vpm MockVenuePizzaModel) Get(id int64) (*Venue, error) {
	return nil, nil
}

func (vpm MockVenuePizzaModel) Update(venue *Venue) error {
	return nil
}

func (vpm MockVenuePizzaModel) Delete(id int64) error {
	return nil
}