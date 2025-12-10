package category

import (
	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/models/category"
)

func (repo *Repository) Create(model *category.Category) (rowID uuid.UUID, err error) {
	err = repo.db.QueryRow(`INSERT INTO categories (`+
		// id
		`name
		,slug
		,position`+
		// ,created_at
		// ,updated_at
		// ,deleted_at
		`) VALUES (
		 $1
		,$2
		,(SELECT COALESCE(MAX(position), 0) + 1 FROM categories WHERE deleted_at IS NULL)`+
		// ,$4
		// ,$5
		// ,$6
		`) RETURNING id`,
		// model.ID,
		model.Name,
		model.Slug,
		// model.Position,
		// model.CreatedAt,
		// model.UpdatedAt,
		// model.DeletedAt,
	).Scan(&rowID)
	return rowID, repo.db.ConvertError(err)
}
