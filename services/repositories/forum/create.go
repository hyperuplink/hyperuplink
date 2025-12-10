package forum

import (
	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/models/forum"
)

func (repo *Repository) Create(model *forum.Forum) (rowID uuid.UUID, err error) {
	err = repo.db.QueryRow(`INSERT INTO forums (
		 name
		,slug
		,position
		,category_id
		,description`+
		// ,created_at
		// ,updated_at
		// ,deleted_at
		`) VALUES (
		,$1
		,$2
		,(SELECT COALESCE(MAX(position), 0) + 1 FROM categories WHERE category_id = $3 AND deleted_at IS NULL)
		,$3
		,$4`+
		// ,$7
		// ,$8
		// ,$9
		`) RETURNING id`,
		model.Name,
		model.Slug,
		// model.Position,
		model.CategoryID,
		model.Description,
		// model.CreatedAt,
		// model.UpdatedAt,
		// model.DeletedAt,
	).Scan(&rowID)
	return rowID, repo.db.ConvertError(err)
}
