package category

import (
	"context"

	"github.com/mrusme/hyperuplink/models/category"
)

func (repo *Repository) Delete(model *category.Category) (err error) {
	tx, cancel, err := repo.db.Tx()
	defer cancel()
	if err != nil {
		return repo.db.ConvertError(err)
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `
		UPDATE categories SET
			updated_at = NOW(),
			deleted_at = NOW()
		WHERE id = $1
		`,
		model.ID,
	)
	if err != nil {
		return repo.db.ConvertError(err)
	}
	_, err = tx.Exec(context.Background(), `
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
	if err != nil {
		return repo.db.ConvertError(err)
	}

	return repo.db.ConvertError(
		tx.Commit(context.Background()),
	)
}
