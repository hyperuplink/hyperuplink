package forum

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
)

func (repo *Repository) Create(ctx context.Context, model *forum.Forum) (rowID uuid.UUID, err error) {
	tx, err := repo.db.Tx(ctx)
	if err != nil {
		return rowID, database.ConvertError(err)
	}
	defer tx.End(ctx)

	err = tx.QueryRow(ctx, `INSERT INTO forums (
		 name
		,slug
		,position
		,category_id
		,description
		,created_at
		,updated_at
	) VALUES (
		 $1
		,$2
		,(SELECT COALESCE(MAX(position), 0) + 1 FROM forums WHERE category_id = $3 AND deleted_at IS NULL)
		,$3
		,$4
		,NOW()
		,NOW()
	) RETURNING id`,
		model.Name,
		model.Slug,
		model.CategoryID,
		model.Description,
	).Scan(&rowID)
	if err != nil {
		return rowID, database.ConvertError(err)
	}

	return rowID, database.ConvertError(tx.Commit(ctx))
}
