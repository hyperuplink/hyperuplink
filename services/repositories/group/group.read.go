package group

import (
	"context"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/group"
)

func (repo *Repository) All(
	ctx context.Context,
	qo common.QueryOptions,
) (model []group.Group, err error) {
	return database.Rows[group.Group](ctx, repo.db, qo.Query(
		`SELECT * FROM groups`,
		common.QueryCapabilities{
			HasDeleted: true,
		}),
	)
}

func (repo *Repository) GetByID(
	ctx context.Context,
	id string,
	qo common.QueryOptions,
) (model *group.Group, err error) {
	return database.Row[group.Group](ctx, repo.db, qo.Query(
		`SELECT * FROM groups WHERE id = $1`,
		common.QueryCapabilities{
			HasDeleted: true,
		}),
		id,
	)
}

func (repo *Repository) AllByIDs(
	ctx context.Context,
	ids []string,
	qo common.QueryOptions,
) (model []group.Group, err error) {
	return database.Rows[group.Group](ctx, repo.db, qo.Query(
		`SELECT * FROM groups WHERE id = ANY($1)`,
		common.QueryCapabilities{
			HasDeleted: true,
		}),
		ids,
	)
}
