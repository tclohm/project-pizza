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

type Venue struct {
	ID int64 `json:"id"`
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
	Address string `json:"address"`
}