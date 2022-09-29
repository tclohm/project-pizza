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

type Venue struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Address string `json:"address"`
}

func ValidateVenue(v *validator.Validator, venue *Venue) {
	v.Check(venue.Name != "", "name", "must be provided")
	v.Check(len(venue.Name) < 500, "name", "must not be more than 500 bytes long")
	v.Check(venue.Address != "", "address", "must be provided")
}

type VenueModel struct {
	DB *sql.DB
}

func (vm VenueModel) Insert(venue *Venue) error {

	query := `SELECT id FROM venues WHERE name = $1 AND address = $2`

	args := []interface{}{
		venue.Name,
		venue.Address,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	exist := vm.DB.QueryRowContext(ctx, query, args...).Scan(&venue.ID)

	if exist != nil && errors.Is(exist, sql.ErrNoRows) {

		query = `
		INSERT INTO venues (
			name, 
			lat,
			lon,
			address
		) VALUES ($1, $2, $3, $4)
		RETURNING id
		`
		// args slices containing values for the placeholder parameters from the venue struct
		args = []interface{}{
			venue.Name, venue.Lat, venue.Lon, venue.Address,
		}

		ctx, cancel = context.WithTimeout(context.Background(), 3 * time.Second)
		defer cancel()

		// passing in the slice and scanning the system generated id
		return vm.DB.QueryRowContext(ctx, query, args...).Scan(&venue.ID)
		
	}


	return exist

	
}

func (vm VenueModel) Get(id int64) (*Venue, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT id, 
		name, 
		lat,
		lon,
		address,
	FROM venues WHERE id = $1
	`

	var venue Venue
	// 3-second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	// release resources associated with context before Get() is returned
	// memory leak avoided
	defer cancel()

	err := vm.DB.QueryRowContext(ctx, query, id).Scan(
		&venue.ID,
		&venue.Name,
		&venue.Lat,
		&venue.Lon,
		&venue.Address,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &venue, nil
}

func (vm VenueModel) Update(venue *Venue) error {
	query := `
		UPDATE venues
		SET name = $1,
		lat = $2,
		lon = $3,
		address = $4,
		WHERE id = $5
		RETURNING id
	`

	args := []interface{}{
		venue.Name,
		venue.Lat,
		venue.Lon,
		venue.Address,
		venue.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()
	// query and scan the new value in
	err := vm.DB.QueryRowContext(ctx, query, args...).Scan(&venue.ID)
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

func (vm VenueModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM venues
		WHERE id = $1`

	result, err := vm.DB.Exec(query, id)
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

func (vm VenueModel) GetAll() ([]*Venue, error) {
	query := `
		SELECT 
		id, 
		name, 
		lat,
		lon,
		address
		FROM venues
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	args := []interface{}{}

	rows, err := vm.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	venues := []*Venue{}

	for rows.Next() {
		var venue Venue

		err := rows.Scan(
			&venue.ID,
			&venue.Name,
			&venue.Lat,
			&venue.Lon,
			&venue.Address,
		)

		if err != nil {
			return nil, err
		}

		venues = append(venues, &venue)

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return venues, nil
}


type MockVenueModel struct {}

func (vm MockVenueModel) Insert(venue *Venue) error {
	return nil
}

func (vm MockVenueModel) Get(id int64) (*Venue, error) {
	return nil, nil
}

func (vm MockVenueModel) Update(venue *Venue) error {
	return nil
}

func (vm MockVenueModel) Delete(id int64) error {
	return nil
}

func (vm MockVenueModel) GetAll() ([]*Venue, error) {
	return nil, nil
}