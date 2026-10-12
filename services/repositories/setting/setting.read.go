package setting

import (
	"context"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
)

func GetByID[T any](
	ctx context.Context,
	repo *Repository,
	id string,
) (model *setting.Setting[T], err error) {
	return database.Row[setting.Setting[T]](ctx, repo.db, `SELECT * FROM settings
		WHERE id = $1
		LIMIT 1`,
		id,
	)
}
