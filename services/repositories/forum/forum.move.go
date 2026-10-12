package forum

import (
	"context"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
)

func (repo *Repository) MoveUp(model *forum.Forum) (err error) {
	_, err = database.Exec(context.Background(), repo.db, `WITH moved_up AS (
		SELECT id, position, category_id
		FROM forums
		WHERE id = $1
			AND deleted_at IS NULL
	)
	UPDATE forums
	SET position = CASE
		WHEN forums.id = moved_up.id THEN moved_up.position - 1
		WHEN forums.position = moved_up.position - 1 AND forums.category_id = moved_up.category_id THEN moved_up.position
		ELSE forums.position
	END
	FROM moved_up
	WHERE (forums.id = moved_up.id OR (forums.position = moved_up.position - 1 AND forums.category_id = moved_up.category_id))
		AND forums.category_id = moved_up.category_id
		AND forums.deleted_at IS NULL`,
		model.ID,
	)
	return err
}

func (repo *Repository) MoveDown(model *forum.Forum) (err error) {
	_, err = database.Exec(context.Background(), repo.db, `WITH moved_down AS (
		SELECT id, position, category_id
		FROM forums
		WHERE id = $1
			AND deleted_at IS NULL
	)
	UPDATE forums
	SET position = CASE
		WHEN forums.id = moved_down.id THEN moved_down.position + 1
		WHEN forums.position = moved_down.position + 1 AND forums.category_id = moved_down.category_id THEN moved_down.position
		ELSE forums.position
	END
	FROM moved_down
	WHERE (forums.id = moved_down.id OR (forums.position = moved_down.position + 1 AND forums.category_id = moved_down.category_id))
		AND forums.category_id = moved_down.category_id
		AND forums.deleted_at IS NULL`,
		model.ID,
	)
	return err
}
