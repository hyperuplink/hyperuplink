package apikey

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
)

func (repo *Repository) TouchLastUsed(ctx context.Context, id uuid.UUID) (err error) {
	_, err = database.Exec(
		ctx,
		repo.db,
		`UPDATE apikeys SET last_used_at = NOW() WHERE id = $1`,
		id,
	)
	return err
}
