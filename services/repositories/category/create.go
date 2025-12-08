package category

import (
	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/models/category"
)

func (repo *Repository) Create(model *category.Category) (rowID uuid.UUID, err error) {
	err = repo.db.QueryRow(`INSERT INTO categories (
		 id
		,name
		,slug`+
		// ,created_at
		// ,updated_at
		// ,deleted_at
		`) VALUES (
		,$1
		,$2
		,$3`+
		// ,$4
		// ,$5
		// ,$6
		`) RETURNING id`,
		model.ID,
		model.Name,
		model.Slug,
		// model.CreatedAt,
		// model.UpdatedAt,
		// model.DeletedAt,
	).Scan(&rowID)
	return rowID, repo.db.ConvertError(err)
}
