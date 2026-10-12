package category

import (
	"context"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
)

func (repo *Repository) Delete(ctx context.Context, model *category.Category) (err error) {
	tx, err := repo.db.Tx(ctx)
	if err != nil {
		return database.ConvertError(err)
	}
	defer tx.End(ctx)

	_, err = database.Exec(ctx, tx, `
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
		return err
	}

	_, err = database.Exec(ctx, tx, `
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
		return err
	}

	_, err = database.Exec(ctx, tx, `
		UPDATE forums SET
			updated_at = NOW(),
			deleted_at = NOW()
		WHERE deleted_at IS NULL
		AND category_id = $1
		`,
		model.ID,
	)
	if err != nil {
		return err
	}

	_, err = database.Exec(ctx, tx, `
		UPDATE categories SET
			updated_at = NOW(),
			deleted_at = NOW()
		WHERE id = $1
		`,
		model.ID,
	)
	if err != nil {
		return err
	}
	_, err = database.Exec(ctx, tx, `
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
		return err
	}

	return database.ConvertError(tx.Commit(ctx))
}
