package forum

import (
	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/models/forum"
)

func (repo *Repository) Create(model *forum.Forum) (rowID uuid.UUID, err error) {
	err = repo.db.QueryRow(`INSERT INTO forums (
		 id
		,name
		,slug
		,category_id`+
		// ,created_at
		// ,updated_at
		// ,deleted_at
		`) VALUES (
		,$1
		,$2
		,$3
		,$4`+
		// ,$5
		// ,$6
		// ,$7
		`) RETURNING id`,
		model.ID,
		model.Name,
		model.Slug,
		model.CategoryID,
		// model.CreatedAt,
		// model.UpdatedAt,
		// model.DeletedAt,
	).Scan(&rowID)
	return rowID, repo.db.ConvertError(err)
}
