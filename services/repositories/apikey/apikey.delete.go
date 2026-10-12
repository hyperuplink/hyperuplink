package apikey

import (
	"context"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/apikey"
)

func (repo *Repository) Delete(ctx context.Context, model *apikey.APIKey) (err error) {
	_, err = database.Exec(
		ctx,
		repo.db,
		`UPDATE apikeys SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1`,
		model.ID,
	)
	return err
}
