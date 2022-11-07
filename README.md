# project-pizza-api

## Different Types of Pizza

### go run ./cmd/api -help

#### -cors-trusted-origins value
    Trusted CORS origins (space separated)
#### -db-ds string
    PostgreSQL DSN (default "postgres://username:password@host/database_name?sslmode=disable")
#### -db-max-idle-conns int
    PostgreSQL max idle connections (default 25)
#### -db-max-idle-time string
    PostgreSQL max connection idle time (default "15m")
#### -db-max-open-conns int
    PostgreSQL max open connections (default 25)
#### -env string
    Environment (development|staging|production) (default "development")
#### -limiter-burst int
    Rate limiter maximum burst (default 100)
#### -limiter-enabled
    Enable rate limiter (default true)
#### -limiter-rps float
    Rate limiter maximum requests per second (default 50)
#### -port int
    API server port (default 4000)

### Example:
-- go run ./cmd/api -cors-trusted-origins="http://localhost:3000 http://localhost:3000/*"

## Models for DB

```
Reviews interface {
    Insert(review *Review) error
    Get(startDate, endDate string) ([]*ReviewWithPizzaName, error)
    Update(review *Review) error
    Delete(id int64) error
    GetAll() ([]*Review, error)
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
    GetPizza(id int64) (*Opinion, error)
    Get(id int64) (*VenuePizza, error) 
    Update(venuePizza *VenuePizza) error
    Delete(id int64) error
    GetAll() ([]*VenuePizzaMixin, error)
}
```
