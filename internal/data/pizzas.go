package data

import (
	"time"
	"database/sql"
	"errors"
	"context"
	"fmt"

	"github.com/tclohm/project-pizza/internal/validator"

	_ "github.com/lib/pq"
)

type Pizza struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Style string `json:"style"`
	Description string `json:"description"`
	Cheesiness float32 `json:"cheesiness"`
	Flavor float32 `json:"flavor"`
	Sauciness float32 `json:"sauciness"`
	Saltiness float32 `json:"saltiness"`
	Charness float32 `json:"charness"`
	CreatedAt time.Time `json:"created_at"`
}

func ValidatePizza(v *validator.Validator, pizza *Pizza) {
	v.Check(pizza.Name != "", "name", "must be provided")
	v.Check(len(pizza.Name) < 500, "name", "must not be more than 500 bytes long")

	v.Check(pizza.Style != "", "style", "must be provided")
	v.Check(len(pizza.Style) < 500, "style", "must not be more than 500 bytes long")

	v.Check(pizza.Description != "", "description", "must be provided")
	v.Check(len(pizza.Description) < 500, "description", "must not be more than 500 bytes long")

	v.Check(pizza.Cheesiness >= 0, "cheesiness", "must be greater than or equal to 0")
	v.Check(pizza.Cheesiness <= 5, "cheesiness", "must be less than or equal to 5")

	v.Check(pizza.Flavor >= 0, "flavor", "must be greater than or equal to 0")
	v.Check(pizza.Flavor <= 5, "flavor", "must be less than or equal to 5")

	v.Check(pizza.Sauciness >= 0, "sauciness", "must be greater than or equal to 0")
	v.Check(pizza.Sauciness <= 5, "sauciness", "must be less than or equal to 5")

	v.Check(pizza.Saltiness >= 0, "saltiness", "must be greater than or equal to 0")
	v.Check(pizza.Saltiness <= 5, "saltiness", "must be less than or equal to 5")

	v.Check(pizza.Charness >= 0, "charness", "must be greater than or equal to 0")
	v.Check(pizza.Charness <= 5, "charness", "must be less than or equal to 5")
}

// wrapper around db connection pool
type PizzaModel struct {
	DB *sql.DB
}

func (pm PizzaModel) Insert(pizza *Pizza) error {
	query := `
	INSERT INTO pizzas (
		name, 
		style, 
		description, 
		cheesiness, 
		flavor, 
		sauciness, 
		saltiness, 
		charness
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id
	`
	// args slices containing values for the placeholder parameters from the pizza struct
	args := []interface{}{
		pizza.Name, pizza.Style, pizza.Description, pizza.Cheesiness, 
		pizza.Flavor, pizza.Sauciness, pizza.Saltiness, pizza.Charness,
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
		style,
		description, 
		cheesiness, 
		flavor, 
		sauciness, 
		saltiness, 
		charness
		FROM pizzas WHERE id = $1
	`

	var pizza Pizza
	// 3-second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	// release resources associated with context before Get() is returned
	// memory leak avoided
	defer cancel()

	err := pm.DB.QueryRowContext(ctx, query, id).Scan(
		&pizza.ID,
		&pizza.Name,
		&pizza.Style,
		&pizza.Description,
		&pizza.Cheesiness,
		&pizza.Flavor,
		&pizza.Sauciness,
		&pizza.Saltiness,
		&pizza.Charness,
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
		UPDATE pizzas
		SET name = $1,
		description = $2, 
		cheesiness = $3, 
		flavor = $4, 
		sauciness = $5, 
		saltiness = $6, 
		charness = $7
		WHERE id = $8
		RETURNING id
	`

	args := []interface{}{
		pizza.Name,
		pizza.Description,
		pizza.Cheesiness,
		pizza.Flavor,
		pizza.Sauciness,
		pizza.Saltiness,
		pizza.Charness,
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

func (pm PizzaModel) GetAll(name string, style string, filters Filters) ([]*Pizza, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT 
		count(*) OVER(),
		id,
		name,
		style,
		description, 
		cheesiness, 
		flavor, 
		sauciness, 
		saltiness, 
		charness
		FROM pizzas
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (LOWER(style) = LOWER($2) OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, 
		filters.sortColumn(), filters.sortDirection()) 

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	args := []interface{}{name, style, filters.limit(), filters.offset()}

	rows, err := pm.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	pizzas := []*Pizza{}

	for rows.Next() {
		var pizza Pizza

		err := rows.Scan(
			&totalRecords,
			&pizza.ID,
			&pizza.Name,
			&pizza.Style,
			&pizza.Description, 
			&pizza.Cheesiness, 
			&pizza.Flavor, 
			&pizza.Sauciness, 
			&pizza.Saltiness, 
			&pizza.Charness,
		)

		if err != nil {
			return nil, Metadata{}, err
		}

		pizzas = append(pizzas, &pizza)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return pizzas, metadata, nil
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

func (pm MockPizzaModel) GetAll(name string, style string, filters Filters) ([]*Pizza, Metadata, error) {
	return nil, Metadata{}, nil
}