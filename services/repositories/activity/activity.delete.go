package activity

import (
	"context"

	"xn--gckvb8fzb.com/glides/services/database"
)

func (repo *Repository) DeleteOlderThanDays(
	ctx context.Context,
	kinds []string,
	days int,
) (deleted int64, err error) {
	if len(kinds) == 0 || days < 1 {
		return 0, nil
	}

	tag, err := database.Exec(ctx, repo.db, `DELETE FROM activities
		WHERE kind = ANY($1)
		AND created_at < NOW() - make_interval(days => $2)`,
		kinds,
		days,
	)
	if err != nil {
		return 0, err
	}

	return tag.RowsAffected(), nil
}
