package group

import (
	"context"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/group"
)

func (repo *Repository) Create(ctx context.Context, model *group.Group) (rowID string, err error) {
	return database.Value[string](ctx, repo.db, `INSERT INTO groups (
		 id
		,name
		,created_at
		,updated_at
		,deleted_at
		) VALUES (
		 $1
		,$2
		,NOW()
		,NOW()
		,NULL
		) RETURNING id`,
		model.ID,
		model.Name,
	)
}
