package forum

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
)

func (repo *Repository) Create(model *forum.Forum) (rowID uuid.UUID, err error) {
	tx, err := repo.db.Tx()
	if err != nil {
		return rowID, repo.db.ConvertError(err)
	}
	defer tx.End()

	err = tx.QueryRow(`INSERT INTO forums (
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
		return rowID, repo.db.ConvertError(err)
	}

	return rowID, repo.db.ConvertError(tx.Commit())
}
