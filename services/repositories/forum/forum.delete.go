package forum

import (
	"context"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
)

func (repo *Repository) Delete(ctx context.Context, model *forum.Forum) (err error) {
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
			SELECT id FROM topics
			WHERE forum_id = $1
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
		AND forum_id = $1
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
		WHERE id = $1
		`,
		model.ID,
	)
	if err != nil {
		return err
	}

	_, err = database.Exec(ctx, tx, `
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
		return err
	}

	return database.ConvertError(tx.Commit(ctx))
}
