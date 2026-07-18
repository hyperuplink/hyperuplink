package category

import (
	"xn--gckvb8fzb.com/hyperuplink/models/category"
)

func (repo *Repository) Delete(model *category.Category) (err error) {
	tx, err := repo.db.Tx()
	if err != nil {
		return repo.db.ConvertError(err)
	}
	defer tx.End()

	_, err = tx.Exec(`
		UPDATE replies SET
			updated_at = NOW(),
			deleted_at = NOW()
		WHERE deleted_at IS NULL
		AND topic_id IN (
			SELECT t.id FROM topics t
			JOIN forums f ON f.id = t.forum_id
			WHERE f.category_id = $1
		)
		`,
		model.ID,
	)
	if err != nil {
		return repo.db.ConvertError(err)
	}

	_, err = tx.Exec(`
		UPDATE topics SET
			updated_at = NOW(),
			deleted_at = NOW()
		WHERE deleted_at IS NULL
		AND forum_id IN (
			SELECT id FROM forums
			WHERE category_id = $1
		)
		`,
		model.ID,
	)
	if err != nil {
		return repo.db.ConvertError(err)
	}

	_, err = tx.Exec(`
		UPDATE forums SET
			updated_at = NOW(),
			deleted_at = NOW()
		WHERE deleted_at IS NULL
		AND category_id = $1
		`,
		model.ID,
	)
	if err != nil {
		return repo.db.ConvertError(err)
	}

	_, err = tx.Exec(`
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
	_, err = tx.Exec(`
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

	return repo.db.ConvertError(tx.Commit())
}
