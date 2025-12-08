package topic

import (
	"github.com/mrusme/hyperuplink/services/database"
)

type Repository struct {
	db *database.Database
}

func New(db *database.Database) (*Repository, error) {
	repo := new(Repository)
	repo.db = db

	return repo, nil
}

func (repo *Repository) Startup() (err error) {
	return nil
}

func (repo *Repository) Shutdown() error {
	return nil
}

