package apikey

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/apikey"
)

func (repo *Repository) GetBySecretHash(
	secretHash string,
	qo common.QueryOptions,
) (model *apikey.APIKey, err error) {
	var rows pgx.Rows
	var mod apikey.APIKey

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM apikeys WHERE secret_hash = $1`,
		common.QueryCapabilities{
			HasDeleted: true,
		}),
		secretHash,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[apikey.APIKey])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetByUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *apikey.APIKey, err error) {
	var rows pgx.Rows
	var mod apikey.APIKey

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM apikeys WHERE id = $1`,
		common.QueryCapabilities{
			HasDeleted: true,
		}),
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[apikey.APIKey])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) AllForUserUUID(
	userID uuid.UUID,
	qo common.QueryOptions,
) (model *[]apikey.APIKey, err error) {
	var rows pgx.Rows
	var mod []apikey.APIKey

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM apikeys WHERE user_id = $1`,
		common.QueryCapabilities{
			HasDeleted: true,
		}),
		userID,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[apikey.APIKey])

	return &mod, repo.db.ConvertError(err)
}
