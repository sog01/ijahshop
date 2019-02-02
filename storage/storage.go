package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Storage is entity of storage package
// Storage include connection to database (currently using mysql)
type Storage struct {
	DB *sqlx.DB
}

// New to create new instance of storage package
func New(dataSource string) (Storage, error) {
	db, err := sqlx.Connect("sqlite3", dataSource)
	if err != nil {
		return Storage{}, err
	}

	return Storage{
		DB: db,
	}, nil
}
