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
	var mod []vreply.VReply

	if total, err = repo.VAllCountForTopicUUID(ctx, id, qo); err != nil {
		return nil, 0, err
	}

	mod, err = database.Rows[vreply.VReply](ctx, repo.db, qo.Query(
		`SELECT * FROM vreplies WHERE topic_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)

	return mod, total, err
}

func (repo *Repository) VAllForTopicID(
	ctx context.Context,
	id string,
	qo common.QueryOptions,
) (model []vreply.VReply, total int64, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, total, database.ConvertError(err)
	}

	return repo.VAllForTopicUUID(ctx, uuID, qo)
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
