package forum

import (
	"github.com/mrusme/hyperuplink/models/forum"
)

func (repo *Repository) Delete(model *forum.Forum) (err error) {
	tx, err := repo.db.Tx()
	if err != nil {
		return repo.db.ConvertError(err)
	}
	defer tx.End()

	_, err = tx.Exec(`
		UPDATE forums SET
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
		WITH ordered_forums AS (
			SELECT id, ROW_NUMBER() OVER (PARTITION BY category_id ORDER BY position) - 1 AS new_position
			FROM forums
			WHERE deleted_at IS NULL
		)
		UPDATE forums
		SET position = ordered_forums.new_position
		FROM ordered_forums
		WHERE forums.id = ordered_forums.id
		`,
	)
	if err != nil {
		return repo.db.ConvertError(err)
	}

	_, err = tx.Exec(`REFRESH MATERIALIZED VIEW vforums`)
	if err != nil {
		return repo.db.ConvertError(err)
	}

	return repo.db.ConvertError(tx.Commit())
}
