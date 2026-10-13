package category

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
)

func (repo *Repository) Create(ctx context.Context, model *category.Category) (rowID uuid.UUID, err error) {
	defer repo.memo.Reset()

	return database.Value[uuid.UUID](ctx, repo.db, `INSERT INTO categories (
		 name
		,slug
		,position
		,created_at
		,updated_at
		) VALUES (
		 $1
		,$2
		,(SELECT COALESCE(MAX(position), 0) + 1 FROM categories WHERE deleted_at IS NULL)
		,NOW()
		,NOW()
		) RETURNING id`,
		model.Name,
		model.Slug,
	)
}
