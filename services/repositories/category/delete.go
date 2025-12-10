package category

import (
	"github.com/mrusme/hyperuplink/models/category"
)

func (repo *Repository) Delete(model *category.Category) (err error) {
	_, err = repo.db.Exec(`
		UPDATE categories SET
			updated_at = NOW(),
			deleted_at = NOW()
		WHERE id = $1
		`,
		model.ID,
	)
	_, err = repo.db.Exec(`
		WITH ordered_categories AS (
			SELECT id, ROW_NUMBER() OVER (ORDER BY position) - 1 AS new_position
			FROM categories
			WHERE deleted_at IS NULL
		)
		UPDATE categories
		SET position = ordered_categories.new_position
		FROM ordered_categories
		WHERE categories.id = ordered_categories.id
		`,
	)
	return repo.db.ConvertError(err)
}
