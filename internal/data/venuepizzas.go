package data

import (
	"time"
	"database/sql"
	"errors"
	"context"
	_ "fmt"
	_ "github.com/lib/pq"

	"github.com/tclohm/project-pizza/internal/validator"
)

type VenuePizza struct {
	ID int64 `json:"id"`
	VenueId int64 `json:"venue_id"`
	PizzaId int64 `json:"pizza_id"`
}

type VenuePizzaMixin struct {
	ID 					int64 		`json:"id"`
	PizzaName 			string 		`json:"pizza_name"`
	PizzaStyle 			string 		`json:"pizza_style"`
	Cheesiness 			float32 	`json:"cheesiness"`
	Flavor				float32		`json:"flavor"`
	Sauciness 			float32 	`json:"sauciness"`
	Saltiness  			float32 	`json:"saltiness"`
	Charness 			float32 	`json:"charness"`
	PizzaImageFilename 	string 		`json:"pizza_image_filename"`
	PizzaImageID		int64 		`json:"pizza_image_id"`
	VenueName 			string 		`json:"venue_name"`
	VenueLat 			float64 	`json:"venue_lat"`
	VenueLon 			float64 	`json:"venue_lon"` 			
	VenueAddress 		string 		`json:"venue_address"`
}

func ValidateVenuePizza(v *validator.Validator, venuepizza *VenuePizza) {
	v.Check(venuepizza.PizzaId > 0, "pizza id", "must be greater than 0")
	v.Check(venuepizza.VenueId > 0, "venue id", "must be greater than 0")
}

type VenuePizzaModel struct {
	DB *sql.DB
}

func (vpm VenuePizzaModel) Insert(venuePizza *VenuePizza) error {
	query := `
	INSERT INTO venuepizzas (
		venue_id, pizza_id
	) VALUES ($1, $2)
	RETURNING id
	`
	// args slices containing values for the placeholder parameters from the venue struct
	args := []interface{}{
		venuePizza.VenueId, venuePizza.PizzaId,
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
		&venuepizza.PizzaId,
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
		venuePizza.ID,
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

func (vpm VenuePizzaModel) GetAll() ([]*VenuePizzaMixin, error) {
	query := `
	select 
	venuepizzas.id,
	pizzas.name as pizza_name,
	pizzas.style as pizza_style,
	pizzas.cheesiness,
	pizzas.flavor,
	pizzas.sauciness,
	pizzas.saltiness,
	pizzas.charness,
	images.filename,
	images.id as pizza_image_id,
	venues.name as venue_name,
	venues.lat,
	venues.lon,
	venues.address
	FROM pizzas 
	JOIN images 
	ON pizzas.image_id = images.id
	JOIN venuepizzas
	ON pizzas.id = venuepizzas.pizza_id
	JOIN venues
	ON venues.id = venuepizzas.venue_id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	args := []interface{}{}

	rows, err := vpm.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	venuepizzas := []*VenuePizzaMixin{}

	for rows.Next() {
		var venuepizzaMixin VenuePizzaMixin

		err := rows.Scan(
			&venuepizzaMixin.ID,
			&venuepizzaMixin.PizzaName,
			&venuepizzaMixin.PizzaStyle,
			&venuepizzaMixin.Cheesiness,
			&venuepizzaMixin.Flavor,
			&venuepizzaMixin.Sauciness,
			&venuepizzaMixin.Saltiness,
			&venuepizzaMixin.Charness,
			&venuepizzaMixin.PizzaImageFilename,
			&venuepizzaMixin.PizzaImageID,
			&venuepizzaMixin.VenueName,
			&venuepizzaMixin.VenueLat,
			&venuepizzaMixin.VenueLon,
			&venuepizzaMixin.VenueAddress,
		)

		if err != nil {
			return nil, err
		}

		venuepizzas = append(venuepizzas, &venuepizzaMixin)

		if err = rows.Err(); err != nil {
			return nil, err
		}

	}
	return venuepizzas, nil
}


type MockVenuePizzaModel struct {}

func (vpm MockVenuePizzaModel) Insert(venuePizza *VenuePizza) error {
	return nil
}

func (vpm MockVenuePizzaModel) Get(id int64) (*VenuePizza, error) {
	return nil, nil
}

func (vpm MockVenuePizzaModel) Update(venuePizza *VenuePizza) error {
	return nil
}

func (vpm MockVenuePizzaModel) Delete(id int64) error {
	return nil
}

func (vpm MockVenuePizzaModel) GetAll() ([]*VenuePizzaMixin, error) {
	return nil, nil
}