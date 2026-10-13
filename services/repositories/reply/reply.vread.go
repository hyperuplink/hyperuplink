package reply

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/vreply"
)

func (repo *Repository) VAllCountForTopicUUID(
	ctx context.Context,
	id uuid.UUID,
	qo common.QueryOptions,
) (total int64, err error) {

	qoc := qo
	qoc.OrderBy = ""
	qoc.Limit = 0
	qoc.Page = 0
	return database.Value[int64](ctx, repo.db, qoc.Query(
		`SELECT COUNT(id) AS total FROM vreplies WHERE topic_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)
}

func (repo *Repository) VAllForTopicUUID(
	ctx context.Context,
	id uuid.UUID,
	qo common.QueryOptions,
) (model []vreply.VReply, total int64, err error) {
	if total, err = repo.VAllCountForTopicUUID(ctx, id, qo); err != nil {
		return nil, 0, err
	}

	qoc := qo
	qoc.OrderBy = ""
	qoc.Limit = 0
	qoc.Page = 0
	query := qoc.Query(
		`SELECT * FROM vreplies WHERE topic_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		},
	) + ` ORDER BY created_at ASC, id ASC`

	if qo.Limit > 0 {
		model, err = database.Rows[vreply.VReply](ctx, repo.db,
			query+` LIMIT $2 OFFSET $3`,
			id,
			qo.Limit,
			qo.Offset(),
		)
	} else {
		model, err = database.Rows[vreply.VReply](ctx, repo.db, query, id)
	}

	return model, total, err
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

func (repo *Repository) VAllWithTopicForAuthorUUID(
	ctx context.Context,
	id uuid.UUID,
	qo common.QueryOptions,
) (model []vreply.VReplyWithTopic, err error) {
	return database.Rows[vreply.VReplyWithTopic](ctx, repo.db, qo.Query(
		`SELECT vr.*,
			t.name AS topic_name,
			t.slug AS topic_slug,
			f.slug AS forum_slug,
			c.slug AS category_slug
		FROM vreplies vr
		LEFT JOIN topics t ON t.id = vr.topic_id
		LEFT JOIN forums f ON f.id = t.forum_id
		LEFT JOIN categories c ON c.id = f.category_id
		WHERE vr.author_id = $1`,
		common.QueryCapabilities{
			Table:      "vr",
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)
}
