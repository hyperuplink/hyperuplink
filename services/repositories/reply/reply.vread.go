package reply

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"xn--gckvb8fzb.com/hyperuplink/models/vreply"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

func (repo *Repository) VAllCountForTopicUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (total int64, err error) {
	var rows pgx.Rows

	qoc := qo
	qoc.OrderBy = ""
	qoc.Limit = 0
	qoc.Page = 0
	rows, err = repo.db.Query(qoc.Query(
		`SELECT COUNT(id) AS total FROM vreplies WHERE topic_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)
	if err != nil {
		return total, repo.db.ConvertError(err)
	}

	var pag map[string]any
	pag, err = pgx.CollectOneRow(rows, pgx.RowToMap)
	if err != nil {
		return total, repo.db.ConvertError(err)
	}
	total = pag["total"].(int64)

	return total, nil
}

func (repo *Repository) VAllForTopicUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *[]vreply.VReply, total int64, err error) {
	var rows pgx.Rows
	var mod []vreply.VReply

	total, err = repo.VAllCountForTopicUUID(id, qo)

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM vreplies WHERE topic_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)
	if err != nil {
		return nil, total, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[vreply.VReply])

	return &mod, total, repo.db.ConvertError(err)
}

func (repo *Repository) VAllForTopicID(
	id string,
	qo common.QueryOptions,
) (model *[]vreply.VReply, total int64, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, total, repo.db.ConvertError(err)
	}

	return repo.VAllForTopicUUID(uuID, qo)
}

func (repo *Repository) VAllWithTopicForAuthorUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *[]vreply.VReplyWithTopic, err error) {
	var rows pgx.Rows
	var mod []vreply.VReplyWithTopic

	rows, err = repo.db.Query(qo.Query(
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
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[vreply.VReplyWithTopic])

	return &mod, repo.db.ConvertError(err)
}
