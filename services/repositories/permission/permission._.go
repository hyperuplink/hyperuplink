package permission

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
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
	ctx := context.Background()

	_, err = repo.GetDefault(ctx)
	if errors.Is(err, errs.ErrNoRows) {
		return repo.Apply(ctx, pgtype.Text{}, pgtype.UUID{}, permission.ReadWrite)
	}

	return err
}

func (repo *Repository) Shutdown() error {
	return nil
}
