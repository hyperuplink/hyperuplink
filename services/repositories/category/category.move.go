package category

import (
	"context"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
)

func (repo *Repository) MoveUp(model *category.Category) (err error) {
	_, err = database.Exec(context.Background(), repo.db, `WITH moved_up AS (
		SELECT id, position
		FROM categories
		WHERE id = $1
			AND deleted_at IS NULL
	)
	UPDATE categories
	SET position = CASE
		WHEN categories.id = moved_up.id THEN moved_up.position - 1
		WHEN categories.position = moved_up.position - 1 THEN moved_up.position
		ELSE categories.position
	END
	FROM moved_up
	WHERE (categories.id = moved_up.id OR (categories.position = moved_up.position - 1))
		AND categories.deleted_at IS NULL`,
		model.ID,
	)
	return err
}

func (repo *Repository) MoveDown(model *category.Category) (err error) {
	_, err = database.Exec(context.Background(), repo.db, `WITH moved_down AS (
		SELECT id, position
		FROM categories
		WHERE id = $1
			AND deleted_at IS NULL
	)
	UPDATE categories
	SET position = CASE
		WHEN categories.id = moved_down.id THEN moved_down.position + 1
		WHEN categories.position = moved_down.position + 1 THEN moved_down.position
		ELSE categories.position
	END
	FROM moved_down
	WHERE (categories.id = moved_down.id OR (categories.position = moved_down.position + 1))
		AND categories.deleted_at IS NULL`,
		model.ID,
	)
	return err
}
