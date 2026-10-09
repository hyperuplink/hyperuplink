package topic

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/topic"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
)

func orderFor(sort topic.Sort) string {
	switch sort {
	case topic.SortNew:
		return "pinned DESC, created_at DESC, id DESC"
	case topic.SortReplies:
		return "pinned DESC, replies DESC, last_activity_at DESC, id DESC"
	default:
		return "pinned DESC, last_activity_at DESC, id DESC"
	}
}

func offsetFor(limit int, page int) int {
	if page < 1 {
		page = 1
	}

	return (page - 1) * limit
}

func (repo *Repository) VAllForForumUUIDSorted(
	id uuid.UUID,
	sort topic.Sort,
	limit int,
	page int,
) (model *[]vtopic.VTopic, total int64, err error) {
	var rows pgx.Rows
	var mod []vtopic.VTopic

	if total, err = repo.VAllCountForForumUUID(id, common.QueryOptions{}); err != nil {
		return nil, 0, err
	}

	rows, err = repo.db.Query(context.Background(),
		`SELECT * FROM vtopics
		WHERE forum_id = $1
		AND spammed_at IS NULL
		AND deleted_at IS NULL
		ORDER BY `+orderFor(sort)+`
		LIMIT $2 OFFSET $3`,
		id,
		limit,
		offsetFor(limit, page),
	)
	if err != nil {
		return nil, 0, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[vtopic.VTopic])

	return &mod, total, repo.db.ConvertError(err)
}

func (repo *Repository) VAllForReadableSlugsSorted(
	slugs []string,
	sort topic.Sort,
	limit int,
	page int,
) (model *[]vtopic.VTopic, total int64, err error) {
	var rows pgx.Rows
	var mod []vtopic.VTopic

	if total, err = repo.VAllCountForReadableSlugs(slugs, common.QueryOptions{}); err != nil {
		return nil, 0, err
	}

	query := `SELECT * FROM vtopics
		WHERE spammed_at IS NULL
		AND deleted_at IS NULL`
	args := []any{limit, offsetFor(limit, page)}
	if slugs != nil {
		query += ` AND category_slug = ANY($3)`
		args = append(args, slugs)
	}
	query += ` ORDER BY ` + orderFor(sort) + ` LIMIT $1 OFFSET $2`

	rows, err = repo.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, 0, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[vtopic.VTopic])

	return &mod, total, repo.db.ConvertError(err)
}
