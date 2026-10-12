package activity

import (
	"context"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/activity"
	"xn--gckvb8fzb.com/hyperuplink/models/vactivity"
)

func (repo *Repository) VAllAdminCount(
	qo common.QueryOptions,
) (total int64, err error) {
	qoc := qo
	qoc.OrderBy = ""
	qoc.Limit = 0
	qoc.Page = 0

	return database.Value[int64](context.Background(), repo.db, qoc.Query(
		`SELECT COUNT(a.id) AS total FROM activities a WHERE a.kind = ANY($1)`,
		common.QueryCapabilities{
			Table:      "a",
			HasDeleted: true,
		}),
		activity.AdminKinds(),
	)
}

func (repo *Repository) VAllAdmin(
	qo common.QueryOptions,
) (model *[]vactivity.VActivity, total int64, err error) {
	var mod []vactivity.VActivity

	if total, err = repo.VAllAdminCount(qo); err != nil {
		return nil, total, err
	}

	mod, err = database.Rows[vactivity.VActivity](context.Background(), repo.db, qo.Query(
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

	return &mod, total, err
}
