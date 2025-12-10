package forum

import (
	"github.com/mrusme/hyperuplink/models/forum"
)

func (repo *Repository) Delete(model *forum.Forum) (err error) {
	_, err = repo.db.Exec(`
	WITH updated_deleted AS (
		UPDATE forums
		SET deleted_at = NOW()
		WHERE id = $1
		RETURNING id, category_id, position
	),
	ordered_forums AS (
		SELECT id, ROW_NUMBER() OVER (PARTITION BY category_id ORDER BY position) - 1 AS new_position
		FROM forums
		WHERE deleted_at IS NULL
	)
	UPDATE forums
	SET position = ordered_forums.new_position
	FROM ordered_forums
	WHERE categories.id = ordered_forums.id;
		`,
		model.ID,
	)
	return repo.db.ConvertError(err)
}

