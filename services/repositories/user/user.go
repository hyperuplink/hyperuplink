package user

import (
	"github.com/mrusme/hyperuplink/services/database"
)

type Repository struct {
	db database.IDatabase
}

func New(db database.IDatabase) (*Repository, error) {
	repo := new(Repository)
	repo.db = db

	return repo, nil
}

func (repo *Repository) Startup() error {
	return nil
}
