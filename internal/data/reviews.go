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

// I would recommend it, I liked it, It was fine, I didn't like it, It wasn't for for me
type Review struct {
	ID 			int64 		`json:"id"`
	Style 		string 		`json:"style"`
	Price 		float32 	`json:"price"`
	Cheesiness 	float32 	`json:"cheesiness"`
	Flavor 		float32 	`json:"flavor"`
	Sauciness 	float32 	`json:"sauciness"`
	Saltiness 	float32 	`json:"saltiness"`
	Charness 	float32 	`json:"charness"`
	Conclusion 	string 		`json:"conclusion"`
	Spiciness 	float32 	`json:"spiciness"`
	CreatedAt 	time.Time 	`json:"created_at"`
	ImageId 	int64 		`json:"image_id"`
}

type ReviewWithPizzaName struct {
	ID 			int64 		`json:"id"`
	PizzaId 	int64 		`json:"pizza_id"`
	Name 		string		`json:"name"`
	Style 		string 		`json:"style"`
	Price 		float32 	`json:"price"`
	Cheesiness 	float32 	`json:"cheesiness"`
	Flavor 		float32 	`json:"flavor"`
	Sauciness 	float32 	`json:"sauciness"`
	Saltiness 	float32 	`json:"saltiness"`
	Charness 	float32 	`json:"charness"`
	Conclusion 	string 		`json:"conclusion"`
	Spiciness 	float32 	`json:"spiciness"`
	CreatedAt 	time.Time 	`json:"created_at"`
	ImageId 	int64 		`json:"image_id"`
}



func ValidateReview(v *validator.Validator, review *Review) {

	v.Check(review.Style != "", "style", "must be provided")
	v.Check(len(review.Style) < 500, "style", "must not be more than 500 bytes long")

	v.Check(review.Price >= 0, "price", "must be above 0")
	v.Check(review.Price <= 500, "style", "must below 500")

	v.Check(review.Cheesiness >= 0, "cheesiness", "must be greater than or equal to 0")
	v.Check(review.Cheesiness <= 5, "cheesiness", "must be less than or equal to 5")

	v.Check(review.Flavor >= 0, "flavor", "must be greater than or equal to 0")
	v.Check(review.Flavor <= 5, "flavor", "must be less than or equal to 5")

	v.Check(review.Sauciness >= 0, "sauciness", "must be greater than or equal to 0")
	v.Check(review.Sauciness <= 5, "sauciness", "must be less than or equal to 5")

	v.Check(review.Saltiness >= 0, "saltiness", "must be greater than or equal to 0")
	v.Check(review.Saltiness <= 5, "saltiness", "must be less than or equal to 5")

	v.Check(review.Charness >= 0, "charness", "must be greater than or equal to 0")
	v.Check(review.Charness <= 5, "charness", "must be less than or equal to 5")

	v.Check(review.Spiciness >= 0, "spiciness", "must be greater than or equal to 0")
	v.Check(review.Spiciness <= 5, "spiciness", "must be less than or equal to 5")

	v.Check(review.Conclusion != "", "conclusion", "must be provided")
	v.Check(len(review.Conclusion) < 500, "conclusion", "must not be more than 500 bytes long")
	v.Check(validator.In(review.Conclusion, "RECOMMENDED", "SATISFIED", "CONTENT", "DISSATISFIED", "STAY AWAY"), "conclusion", "must be the provided options")
}

type ReviewModel struct {
	DB *sql.DB
}

func (rm ReviewModel) Insert(review *Review) error {
	query := `
	INSERT INTO reviews (
		style,
		price,
		cheesiness, 
		flavor, 
		sauciness, 
		saltiness, 
		charness,
		spiciness,
		conclusion,
		image_id
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id
	`
	// args slices containing values for the placeholder parameters from the review struct
	args := []interface{}{
		review.Style, 
		review.Price, 
		review.Cheesiness, 
		review.Flavor, 
		review.Sauciness, 
		review.Saltiness, 
		review.Charness,
		review.Spiciness,
		review.Conclusion,
		review.ImageId,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	// passing in the slice and scanning the system generated id
	return rm.DB.QueryRowContext(ctx, query, args...).Scan(&review.ID)
}

func (rm ReviewModel) Get(startDate, endDate string) ([]*ReviewWithPizzaName, error) {
	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return nil, err
	}

	endYear, endMonth, endDay := end.Date()

	if start.Before(time.Date(2022, 10, 01, 0, 0, 0, 0, time.UTC)) {
		fmt.Println("error", start)
		return nil, ErrRecordNotFound
	}

	if start.After(time.Date(endYear, endMonth, endDay, 0, 0, 0, 0, time.UTC)) {
		fmt.Println("error start after end", end)
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT 
		reviews.id,
		pizzas.id,
		name,
		style,
		price,
		cheesiness,
		flavor,
		sauciness,
		saltiness,
		charness,
		spiciness,
		conclusion,
		image_id
	FROM reviews 
	JOIN pizzas ON reviews.id = pizzas.review_id
	WHERE created_at BETWEEN $1 and $2
	`

	args := []interface{}{
		start,
		end,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	// release resources associated with context before Get() is returned
	// memory leak avoided
	defer cancel()

	rows, err := rm.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	reviews := []*ReviewWithPizzaName{}

	for rows.Next() {
		var review ReviewWithPizzaName

		err = rows.Scan(
			&review.ID,
			&review.PizzaId,
			&review.Name,
			&review.Style,
			&review.Price,
			&review.Cheesiness,
			&review.Flavor,
			&review.Sauciness,
			&review.Saltiness,
			&review.Charness,
			&review.Spiciness,
			&review.Conclusion,
			&review.ImageId,
		)

		if err != nil {
			return nil, err
		}

		reviews = append(reviews, &review)

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return reviews, nil
}

// MARK: -- UPDATE
func (rm ReviewModel) Update(review *Review) error {
	query := `
	UPDATE reviews
		SET
		style = $1,
		price = $2
		cheesiness = $3, 
		flavor = $4, 
		sauciness = $5, 
		saltiness = $6, 
		charness = $7,
		spiciness = $8,
		conclusion = $9,
		image_id = $10,
	WHERE id = $11
	RETURNING id
	`

	args := []interface{}{
		review.Style,
		review.Price,
		review.Cheesiness,
		review.Flavor,
		review.Sauciness,
		review.Saltiness,
		review.Charness,
		review.Spiciness,
		review.Conclusion,
		review.ImageId,
		review.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()
	// query and scan the new value in
	err := rm.DB.QueryRowContext(ctx, query, args...).Scan(&review.ID)
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

func (rm ReviewModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM 
		reviews
	WHERE id = $1`

	result, err := rm.DB.Exec(query, id)
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

func (rm ReviewModel) GetAll() ([]*Review, error) {
	query := `
		SELECT
			id, 
			style,
			price,
			cheesiness, 
			flavor, 
			sauciness, 
			saltiness, 
			charness,
			spiciness,
			conclusion,
			image_id,
			created_at
		FROM reviews`

	ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
	defer cancel()

	args := []interface{}{}

	rows, err := rm.DB.QueryContext(ctx, query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	reviews := []*Review{}

	for rows.Next() {
		var review Review

		err := rows.Scan(
			&review.ID,
			&review.Style,
			&review.Price,
			&review.Cheesiness,
			&review.Flavor,
			&review.Sauciness,
			&review.Saltiness,
			&review.Charness,
			&review.Spiciness,
			&review.Conclusion,
			&review.ImageId,
			&review.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		reviews = append(reviews, &review)

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return reviews, nil
}


type MockReviewModel struct {}

func (rm MockReviewModel) Insert(review *Review) error {
	return nil
}

func (rm MockReviewModel) Get(startDate, endDate string) ([]*ReviewWithPizzaName, error) {
	return nil, nil
}

func (rm MockReviewModel) Update(review *Review) error {
	return nil
}

func (rm MockReviewModel) Delete(id int64) error {
	return nil
}

func (rm MockReviewModel) GetAll() ([]*Review, error) {
	return nil, nil
}