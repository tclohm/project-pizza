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

type Image struct {
	ID int64 `json:"id"`
	Filename string `json:"filename"`
	ContentType string `json:"content_type"`
	Location string `json:"location"`
	CreatedAt time.Time `json:"created_at"`
}

func ValidateImage(v *validator.Validator, image *Image) {
	v.Check(image.Filename != "", "name", "must be provided")
	v.Check(image.ContentType == "image/png" || image.ContentType == "image/jpeg", "type", "must be either a jpeg or png")
}

type ImageModel struct {
	DB *sql.DB
}

func (im ImageModel) Insert(image *Image) error {
	query := `
	INSERT INTO images (
		filename,
		content_type,
		location
	)
	VALUES($1, $2, $3)
	RETURNING id
	`

	args := []interface{}{
		image.Filename, image.ContentType, image.Location,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	return im.DB.QueryRowContext(ctx, query, args...).Scan(&image.ID)
}

func (im ImageModel) Get(id int64) (*Image, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, 
		filename,
		content_type,
		location
		FROM pizzas WHERE id = $1
	`

	var image Image
	// 3-second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	// release resources associated with context before Get() is returned
	// memory leak avoided
	defer cancel()

	err := im.DB.QueryRowContext(ctx, query, id).Scan(
		&image.Filename,
		&image.ContentType,
		&image.Location,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &image, nil
}

func (im ImageModel) Update(image *Image) error {
	query := `
		UPDATE images
		SET filename = $1,
		content_type = $2, 
		location = $3, 
		WHERE id = $4
		RETURNING id
	`

	args := []interface{}{
		image.Filename,
		image.ContentType,
		image.Location,
		image.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()
	// query and scan the new value in
	err := im.DB.QueryRowContext(ctx, query, args...).Scan(&image.ID)
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

func (im ImageModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM images
		WHERE id = $1`

	result, err := im.DB.Exec(query, id)
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

type MockImageModel struct {}

func (pm MockImageModel) Insert(image *Image) error {
	return nil
}

func (pm MockImageModel) Get(id int64) (*Image, error) {
	return nil, nil
}

func (pm MockImageModel) Update(image *Image) error {
	return nil
}

func (pm MockImageModel) Delete(id int64) error {
	return nil
}