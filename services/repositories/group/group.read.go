package group

import (
	"context"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/group"
)

func (repo *Repository) All(
	qo common.QueryOptions,
) (model *[]group.Group, err error) {
	var mod []group.Group

	mod, err = database.Rows[group.Group](context.Background(), repo.db, qo.Query(
		`SELECT * FROM groups`,
		common.QueryCapabilities{
			HasDeleted: true,
		}),
	)

	return &mod, err
}

func (repo *Repository) GetByID(
	id string,
	qo common.QueryOptions,
) (model *group.Group, err error) {
	return database.Row[group.Group](context.Background(), repo.db, qo.Query(
		`SELECT * FROM groups WHERE id = $1`,
		common.QueryCapabilities{
			HasDeleted: true,
		}),
		id,
	)
}

func (repo *Repository) AllByIDs(
	ids []string,
	qo common.QueryOptions,
) (model *[]group.Group, err error) {
	var mod []group.Group

	mod, err = database.Rows[group.Group](context.Background(), repo.db, qo.Query(
		`SELECT * FROM groups WHERE id = ANY($1)`,
		common.QueryCapabilities{
			HasDeleted: true,
		}),
		ids,
	)

	return &mod, err
}
