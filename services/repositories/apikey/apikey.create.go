package apikey

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/apikey"
)

func (repo *Repository) Create(ctx context.Context, model *apikey.APIKey) (rowID uuid.UUID, err error) {
	return database.Value[uuid.UUID](ctx, repo.db, `INSERT INTO apikeys (
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
	)
}
