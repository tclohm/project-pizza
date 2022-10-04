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

type PizzaReviewed struct {
	VenueId 			int64 		`json:"venue_id"`
	PizzaName 			string 		`json:"pizza_name"`
	Opinions 			[]*Opinion 	`json:"opinions"`
}

type Opinion struct {
	VenueId 			int64 		`json:"venue_id"`
	PizzaId 			int64 		`json:"pizza_id"`
	PizzaName 			string 		`json:"pizza_name"`
	PizzaStyle 			string 		`json:"pizza_style"`
	PizzaPrice			float32 	`json:"price"`
	Cheesiness 			float32 	`json:"cheesiness"`
	Flavor				float32		`json:"flavor"`
	Sauciness 			float32 	`json:"sauciness"`
	Saltiness  			float32 	`json:"saltiness"`
	Charness 			float32 	`json:"charness"`
	Spiciness 			float32		`json:"spiciness"`
	Conclusion 			string 		`json:"conclusion"`
	PizzaImageFilename 	string 		`json:"pizza_image_filename"`
	PizzaImageID		int64 		`json:"pizza_image_id"`
	PizzaLocation		string 		`json:"pizza_image_location"`
	CreatedAt			time.Time 	`json:"created_at"`
}


type VenuePizzaMixin struct {
	VenueId 			int64 			 `json:"venue_id"`
	VenueName 			string 			 `json:"venue_name"`
	Lat 				float64 		 `json:"lat"`
	Lon 				float64 		 `json:"lon"` 			
	VenueAddress 		string 			 `json:"venue_address"`
	Pizzas  			[]*PizzaReviewed `json:"pizzas"`
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

	// type VenuePizzaMixin struct {
	// 	VenueId 			int64 			 `json:"venue_id"`
	// 	VenueName 			string 			 `json:"venue_name"`
	// 	Lat 				float64 		 `json:"lat"`
	// 	Lon 				float64 		 `json:"lon"` 			
	// 	VenueAddress 		string 			 `json:"venue_address"`
	// 	Pizzas  			[]*PizzaReviewed `json:"pizzas"`
	// }

	venue_query := `
		SELECT
			venues.id as venue_id,
			venues.name as venue_name,
			venues.lat,
			venues.lon,
			venues.address
		FROM venuepizzas
		JOIN venues
		ON venues.id = venuepizzas.venue_id
		GROUP BY venues.id
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	venue_args := []interface{}{}

	venue_rows, err := vpm.DB.QueryContext(ctx, venue_query, venue_args...)

	if err != nil {
		return nil, err
	}

	defer venue_rows.Close()

	venuepizzas := []*VenuePizzaMixin{}

	for venue_rows.Next() {
		var venuepizzaMixin VenuePizzaMixin
		var pizzas = []*PizzaReviewed{}

		err := venue_rows.Scan(
			&venuepizzaMixin.VenueId,
			&venuepizzaMixin.VenueName,
			&venuepizzaMixin.Lat,
			&venuepizzaMixin.Lon,
			&venuepizzaMixin.VenueAddress,
		)

		if err != nil {
			return nil, err
		}

		// type PizzaReviewed struct {
			// VenueId 			int64 		`json:"venue_id"`
			// PizzaName 		string 		`json:"pizza_name"`
			// Opinions 		[]*Opinion 	`json:"opinions"`
		// }

		pizza_reviewed_for_venue_query := `
		SELECT
			DISTINCT venuepizzas.venue_id,
			pizzas.name as pizza_name
			FROM venuepizzas
			JOIN pizzas
			ON venuepizzas.pizza_id = pizzas.id
		`

		pizza_args := []interface{}{}

		pizza_rows, err := vpm.DB.QueryContext(ctx, pizza_reviewed_for_venue_query, pizza_args...)

		if err != nil {
			return nil, err
		}

		defer pizza_rows.Close()

		for pizza_rows.Next() {
			var pizzaReviewed PizzaReviewed
			var opinions = []*Opinion{}

			err := pizza_rows.Scan(
				&pizzaReviewed.VenueId,
				&pizzaReviewed.PizzaName,
			)

			if err != nil {
				return nil, err
			}

			// this is where we look for reviews and add them

			opinion_query := `
			select 
				venues.id as venue_id,
				pizzas.id as pizza_id,
				pizzas.name as pizza_name, 
				reviews.style as pizza_style,
				reviews.price,
				reviews.cheesiness,
				reviews.flavor,
				reviews.sauciness,
				reviews.saltiness,
				reviews.charness,
				reviews.spiciness,
				reviews.conclusion,
				images.filename as pizza_image_filename,
				images.id as pizza_image_id,
				images.location as pizza_image_location,
				reviews.created_at
			FROM venues
			JOIN venuepizzas
			ON venues.id = venuepizzas.venue_id
			JOIN pizzas
			ON pizzas.id = venuepizzas.pizza_id
			JOIN reviews
			ON pizzas.review_id = reviews.id
			JOIN images
			ON reviews.image_id = images.id
			ORDER BY reviews.created_at DESC
			`

			opinion_args := []interface{}{}

			opinion_rows, err := vpm.DB.QueryContext(ctx, opinion_query, opinion_args...)

			if err != nil {
				return nil, err
			}

			defer opinion_rows.Close()

			for opinion_rows.Next() {
				var opinion Opinion

				err := opinion_rows.Scan(
					&opinion.VenueId,
					&opinion.PizzaId,
					&opinion.PizzaName,
					&opinion.PizzaStyle, 
					&opinion.PizzaPrice,			
					&opinion.Cheesiness, 			
					&opinion.Flavor,				
					&opinion.Sauciness, 			
					&opinion.Saltiness,  			
					&opinion.Charness, 			
					&opinion.Spiciness,
					&opinion.Conclusion,			
					&opinion.PizzaImageFilename, 	
					&opinion.PizzaImageID,		
					&opinion.PizzaLocation,	
					&opinion.CreatedAt,				
				)

				if err != nil {
					return nil, err
				}

				if opinion.VenueId == pizzaReviewed.VenueId && opinion.PizzaName == pizzaReviewed.PizzaName {
					opinions = append(opinions, &opinion)
				}

				if err = opinion_rows.Err(); err != nil {
					return nil, err
				}
			}

			pizzaReviewed.Opinions = opinions

			if pizzaReviewed.VenueId == venuepizzaMixin.VenueId {
				pizzas = append(pizzas, &pizzaReviewed)
			}

			if err = pizza_rows.Err(); err != nil {
				return nil, err
			}
		}

		venuepizzaMixin.Pizzas = pizzas

		venuepizzas = append(venuepizzas, &venuepizzaMixin)

		if err = venue_rows.Err(); err != nil {
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