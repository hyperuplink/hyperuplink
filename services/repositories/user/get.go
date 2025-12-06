package user

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mrusme/hyperuplink/models/user"
)

func (repo *Repository) GetByID(
	id string,
) (model *user.User, err error) {
	var uuID uuid.UUID
	var rows pgx.Rows
	var usr user.User

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, repo.db.ConvertError(err)
	}

	rows, err = repo.db.Query(`SELECT * FROM users
	WHERE id = $1
	AND banned_at IS NULL
	AND deleted_at IS NULL
	LIMIT 1`,
		uuID)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	usr, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])

	return &usr, repo.db.ConvertError(err)
}

func (repo *Repository) GetByUsername(
	username string,
) (model *user.User, err error) {
	var rows pgx.Rows
	var usr user.User

	rows, err = repo.db.Query(`SELECT * FROM users
	WHERE username = $1
	AND banned_at IS NULL
	AND deleted_at IS NULL
	LIMIT 1`,
		username)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}
	defer rows.Close()

	usr, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])

	return &usr, repo.db.ConvertError(err)
}

func (repo *Repository) GetByEmail(
	email string,
) (model *user.User, err error) {
	var rows pgx.Rows
	var usr user.User

	rows, err = repo.db.Query(`SELECT * FROM users
	WHERE email = $1
	AND banned_at IS NULL
	AND deleted_at IS NULL
	LIMIT 1`,
		email)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}
	defer rows.Close()

	usr, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])

	return &usr, repo.db.ConvertError(err)
}
