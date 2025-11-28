package repository

import (
	"log"

	"exmple.com/database-query-server/internal/database"
)

type Repository struct {
	Postgress database.ClientInterface
}

// New returns a pointer to a new Repository instance or an error.
func New() (*Repository, error) {
	pg, err := database.NewPostgressClient()
	if err != nil {
		log.Fatal(err)
	}

	return &Repository{Postgress: pg}, err
}
