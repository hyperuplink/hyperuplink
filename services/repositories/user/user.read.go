package user

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

func (repo *Repository) All(
	qo common.QueryOptions,
) (model *[]user.User, err error) {
	var rows pgx.Rows
	var mod []user.User

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM users`,
		common.QueryCapabilities{
			HasBanned:  true,
			HasDeleted: true,
		}),
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[user.User])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetByUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *user.User, err error) {
	var rows pgx.Rows
	var mod user.User

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM users WHERE id = $1`,
		common.QueryCapabilities{
			HasBanned:  true,
			HasDeleted: true,
		}),
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetByID(
	id string,
	qo common.QueryOptions,
) (model *user.User, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, repo.db.ConvertError(err)
	}

	return repo.GetByUUID(uuID, qo)
}

func (repo *Repository) GetByUsername(
	username string,
	qo common.QueryOptions,
) (model *user.User, err error) {
	var rows pgx.Rows
	var mod user.User

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM users WHERE username = $1`,
		common.QueryCapabilities{
			HasBanned:  true,
			HasDeleted: true,
		}),
		username)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}
	defer rows.Close()

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetByEmail(
	email string,
	qo common.QueryOptions,
) (model *user.User, err error) {
	var rows pgx.Rows
	var mod user.User

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM users WHERE email = $1`,
		common.QueryCapabilities{
			HasBanned:  true,
			HasDeleted: true,
		}),
		email)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}
	defer rows.Close()

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[user.User])

	return &mod, repo.db.ConvertError(err)
}
