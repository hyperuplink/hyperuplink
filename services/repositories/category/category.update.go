package category

import (
	"github.com/mrusme/hyperuplink/models/category"
)

func (repo *Repository) Update(model *category.Category) (err error) {
	_, err = repo.db.Exec(`UPDATE categories SET
		 name =                        $1
		,slug =                        $2
		,updated_at =                  NOW()
		WHERE id =                     $3`,
		model.Name,
		model.Slug,
		// WHERE id =
		model.ID,
	)
	return repo.db.ConvertError(err)
}
