package reply

import (
	"context"

	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/vreply"
)

func (repo *Repository) VAllForTopicUUIDOrdered(
	ctx context.Context,
	id uuid.UUID,
	limit int,
	page int,
) (model []vreply.VReply, total int64, err error) {
	var mod []vreply.VReply

	if total, err = repo.VAllCountForTopicUUID(ctx, id, common.QueryOptions{}); err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}

	mod, err = database.Rows[vreply.VReply](ctx, repo.db,
		`SELECT * FROM vreplies
		WHERE topic_id = $1
		AND spammed_at IS NULL
		AND deleted_at IS NULL
		ORDER BY created_at ASC, id ASC
		LIMIT $2 OFFSET $3`,
		id,
		limit,
		(page-1)*limit,
	)

	return mod, total, err
}

func (repo *Repository) VGetByUUID(
	ctx context.Context,
	id uuid.UUID,
	qo common.QueryOptions,
) (model *vreply.VReply, err error) {
	return database.Row[vreply.VReply](ctx, repo.db, qo.Query(
		`SELECT * FROM vreplies WHERE id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)
}

func (repo *Repository) VGetByShortID(
	ctx context.Context,
	shortID string,
	qo common.QueryOptions,
) (model *vreply.VReply, err error) {
	return database.Row[vreply.VReply](ctx, repo.db, qo.Query(
		`SELECT * FROM vreplies WHERE short_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		shortID,
	)
}

func (repo *Repository) PositionForUUID(ctx context.Context, id uuid.UUID) (position int64, err error) {
	err = repo.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM replies r2, replies r
		WHERE r.id = $1
		AND r2.topic_id = r.topic_id
		AND r2.deleted_at IS NULL
		AND r2.spammed_at IS NULL
		AND (r2.created_at < r.created_at
			OR (r2.created_at = r.created_at AND r2.id <= r.id))`,
		id,
	).Scan(&position)

	return position, database.ConvertError(err)
}
