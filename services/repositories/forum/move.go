package forum

import (
	"github.com/mrusme/hyperuplink/models/forum"
)

func (repo *Repository) MoveUp(model *forum.Forum) (err error) {
	_, err = repo.db.Exec(`WITH moved_up AS (
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
	if err != nil {
		// TODO: Does it make sense to run this as a transaction?
		_, err = repo.db.Exec(`REFRESH MATERIALIZED VIEW vforums`)
	}
	return repo.db.ConvertError(err)
}

func (repo *Repository) MoveDown(model *forum.Forum) (err error) {
	_, err = repo.db.Exec(`WITH moved_down AS (
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
	if err != nil {
		// TODO: Does it make sense to run this as a transaction?
		_, err = repo.db.Exec(`REFRESH MATERIALIZED VIEW vforums`)
	}
	return repo.db.ConvertError(err)
}
