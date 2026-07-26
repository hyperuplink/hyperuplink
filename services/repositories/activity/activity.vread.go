package activity

import (
	"github.com/jackc/pgx/v5"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/activity"
	"xn--gckvb8fzb.com/hyperuplink/models/vactivity"
)

func (repo *Repository) VAllAdminCount(
	qo common.QueryOptions,
) (total int64, err error) {
	var rows pgx.Rows

	qoc := qo
	qoc.OrderBy = ""
	qoc.Limit = 0
	qoc.Page = 0

	rows, err = repo.db.Query(qoc.Query(
		`SELECT COUNT(a.id) AS total FROM activities a WHERE a.kind = ANY($1)`,
		common.QueryCapabilities{
			Table:      "a",
			HasDeleted: true,
		}),
		activity.AdminKinds(),
	)
	if err != nil {
		return total, repo.db.ConvertError(err)
	}

	var pag map[string]any
	if pag, err = pgx.CollectOneRow(rows, pgx.RowToMap); err != nil {
		return total, repo.db.ConvertError(err)
	}
	total = pag["total"].(int64)

	return total, nil
}

func (repo *Repository) VAllAdmin(
	qo common.QueryOptions,
) (model *[]vactivity.VActivity, total int64, err error) {
	var rows pgx.Rows
	var mod []vactivity.VActivity

	if total, err = repo.VAllAdminCount(qo); err != nil {
		return nil, total, err
	}

	rows, err = repo.db.Query(qo.Query(
		`SELECT
			 a.*
			,COALESCE(u.username, '') AS actor_username
		FROM activities a
		LEFT JOIN users u ON u.id = a.actor_id
		WHERE a.kind = ANY($1)`,
		common.QueryCapabilities{
			Table:      "a",
			HasDeleted: true,
		}),
		activity.AdminKinds(),
	)
	if err != nil {
		return nil, total, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[vactivity.VActivity])

	return &mod, total, repo.db.ConvertError(err)
}
