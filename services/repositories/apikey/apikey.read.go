package apikey

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/apikey"
)

func (repo *Repository) GetBySecretHash(
	secretHash string,
	qo common.QueryOptions,
) (model *apikey.APIKey, err error) {
	return database.Row[apikey.APIKey](context.Background(), repo.db, qo.Query(
		`SELECT * FROM apikeys WHERE secret_hash = $1`,
		common.QueryCapabilities{
			HasDeleted: true,
		}),
		secretHash,
	)
}

func (repo *Repository) GetByUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *apikey.APIKey, err error) {
	return database.Row[apikey.APIKey](context.Background(), repo.db, qo.Query(
		`SELECT * FROM apikeys WHERE id = $1`,
		common.QueryCapabilities{
			HasDeleted: true,
		}),
		id,
	)
}

func (repo *Repository) AllForUserUUID(
	userID uuid.UUID,
	qo common.QueryOptions,
) (model *[]apikey.APIKey, err error) {
	var mod []apikey.APIKey

	mod, err = database.Rows[apikey.APIKey](context.Background(), repo.db, qo.Query(
		`SELECT * FROM apikeys WHERE user_id = $1`,
		common.QueryCapabilities{
			HasDeleted: true,
		}),
		userID,
	)

	return &mod, err
}
