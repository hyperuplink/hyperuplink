package group

import (
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/memo"
)

type Repository struct {
	db   *database.Database
	memo *memo.Memo
}

func New(db *database.Database, board *memo.Memo) (*Repository, error) {
	repo := new(Repository)
	repo.db = db
	repo.memo = board

	return repo, nil
}

func (repo *Repository) Startup() (err error) {
	return nil
}

func (repo *Repository) Shutdown() error {
	return nil
}
