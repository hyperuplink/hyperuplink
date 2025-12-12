package reply

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mrusme/hyperuplink/models/vreply"
	"github.com/mrusme/hyperuplink/services/repositories/common"
)

func (repo *Repository) VAllForTopicUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *[]vreply.VReply, total int64, err error) {
	var rows pgx.Rows
	var mod []vreply.VReply

	qoc := qo
	qoc.OrderBy = ""
	qoc.Limit = 0
	qoc.Page = 0
	rows, err = repo.db.Query(qoc.Query(
		`SELECT COUNT(id) AS total_posts FROM vreplies WHERE topic_id = $1`,
		common.QueryCapabilities{
			HasSpammed: true,
			HasDeleted: true,
		}),
		id,
	)
	if err != nil {
		return nil, total, repo.db.ConvertError(err)
	}

	var pag map[string]any
	pag, err = pgx.CollectOneRow(rows, pgx.RowToMap)
	if err != nil {
		return nil, total, repo.db.ConvertError(err)
	}
	total = pag["total_posts"].(int64)

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
