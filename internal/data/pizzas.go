package data

import (
	"time"
	"database/sql"
	"errors"
	"context"
	_ "fmt"

	"github.com/tclohm/project-pizza/internal/validator"

	_ "github.com/lib/pq"
)

type Pizza struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	ReviewId int64 `json:"review_id"`
}

func ValidatePizza(v *validator.Validator, pizza *Pizza) {
	v.Check(pizza.Name != "", "name", "must be provided")
	v.Check(len(pizza.Name) < 500, "name", "must not be more than 500 bytes long")
}

// wrapper around db connection pool
type PizzaModel struct {
	DB *sql.DB
}

func (pm PizzaModel) Insert(pizza *Pizza) error {
	query := `
	INSERT INTO pizzas (
		name, 
		review_id
	) VALUES ($1, $2)
	RETURNING id
	`
	// args slices containing values for the placeholder parameters from the pizza struct
	args := []interface{}{
		pizza.Name, 
		pizza.ReviewId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	// passing in the slice and scanning the system generated id
	return pm.DB.QueryRowContext(ctx, query, args...).Scan(&pizza.ID)
}

func (pm PizzaModel) Get(id int64) (*Pizza, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, 
		name,
		review_id
	FROM pizzas 
	WHERE id = $1
	`

	var pizza Pizza
	// 3-second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	// release resources associated with context before Get() is returned
	// memory leak avoided
	defer cancel()

	err := pm.DB.QueryRowContext(ctx, query, id).Scan(
		&pizza.ID,
		&pizza.ReviewId,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &pizza, nil
}

func (pm PizzaModel) Update(pizza *Pizza) error {
	query := `
	UPDATE pizzas SET 
		name = $1,
		review_id = $2
	WHERE id = $3
	RETURNING id
	`

	args := []interface{}{
		pizza.Name,
		pizza.ReviewId,
		pizza.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()
	// query and scan the new value in
	err := pm.DB.QueryRowContext(ctx, query, args...).Scan(&pizza.ID)
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

func (pm PizzaModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM pizzas
	WHERE id = $1`

	result, err := pm.DB.Exec(query, id)
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

func (pm PizzaModel) GetAll() ([]*Pizza, error) {
	query := `
	SELECT 
		id,
		name,
		review_id,
	FROM pizzas
	` 

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	args := []interface{}{}

	rows, err := pm.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	pizzas := []*Pizza{}

	for rows.Next() {
		var pizza Pizza

		err := rows.Scan(
			&pizza.ID,
			&pizza.Name,
			&pizza.ReviewId, 
		)

		if err != nil {
			return nil, err
		}

		pizzas = append(pizzas, &pizza)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return pizzas, nil
}


type MockPizzaModel struct {}

func (pm MockPizzaModel) Insert(pizza *Pizza) error {
	return nil
}

func (pm MockPizzaModel) Get(id int64) (*Pizza, error) {
	return nil, nil
}

func (pm MockPizzaModel) Update(pizza *Pizza) error {
	return nil
}

func (pm MockPizzaModel) Delete(id int64) error {
	return nil
}

func (pm MockPizzaModel) GetAll() ([]*Pizza, error) {
	return nil, nil
}