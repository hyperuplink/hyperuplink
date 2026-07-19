package apikey

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/models/apikey"
)

func (repo *Repository) Create(model *apikey.APIKey) (rowID uuid.UUID, err error) {
	err = repo.db.QueryRow(`INSERT INTO apikeys (
		 user_id
		,name
		,secret_hash
		,created_at
		,updated_at
		) VALUES (
		 $1
		,$2
		,$3
		,NOW()
		,NOW()
		) RETURNING id`,
		model.UserID,
		model.Name,
		model.SecretHash,
	).Scan(&rowID)
	return rowID, repo.db.ConvertError(err)
}
