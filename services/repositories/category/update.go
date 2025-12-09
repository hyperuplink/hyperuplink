package category

import (
	"github.com/mrusme/hyperuplink/models/category"
)

func (repo *Repository) Update(model *category.Category) (err error) {
	_, err = repo.db.Exec(`UPDATE categories SET
		 name =                        $1
		,slug =                        $2`+
		// ,created_at =                  $3
		`,updated_at =                  $3
		,deleted_at =                  $4
		WHERE id =                     $5`,
		model.Name,
		model.Slug,
		// model.CreatedAt,
		model.UpdatedAt,
		model.DeletedAt,
		// WHERE id =
		model.ID,
	)
	return repo.db.ConvertError(err)
}
