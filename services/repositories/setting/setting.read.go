package setting

import (
	"github.com/jackc/pgx/v5"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
)

func GetByID[T any](
	repo *Repository,
	id string,
) (model *setting.Setting[T], err error) {
	var rows pgx.Rows
	var mod setting.Setting[T]

	rows, err = repo.db.Query(`SELECT * FROM settings
		WHERE id = $1
		LIMIT 1`,
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[setting.Setting[T]])

	return &mod, repo.db.ConvertError(err)
}
