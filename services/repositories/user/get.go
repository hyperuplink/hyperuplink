package user

import (
	"github.com/jackc/pgx/v5"
	"github.com/mrusme/hyperuplink/models/user"
)

func (repo *Repository) GetByUsername(username string) (model *user.User, err error) {
	var rows pgx.Rows
	var usr user.User
	model = new(user.User)
	rows, err = repo.db.Query(`SELECT * FROM users
	WHERE username = $1
	AND banned_at IS NULL
	AND deleted_at IS NULL
	LIMIT 1`,
		username)

	usr, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])

	return &usr, err
}
