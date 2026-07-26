package group

import (
	"github.com/jackc/pgx/v5"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/models/group"
)

func (repo *Repository) All(
	qo common.QueryOptions,
) (model *[]group.Group, err error) {
	var rows pgx.Rows
	var mod []group.Group

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM groups`,
		common.QueryCapabilities{
			HasDeleted: true,
		}),
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[group.Group])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetByID(
	id string,
	qo common.QueryOptions,
) (model *group.Group, err error) {
	var rows pgx.Rows
	var mod group.Group

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM groups WHERE id = $1`,
		common.QueryCapabilities{
			HasDeleted: true,
		}),
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[group.Group])

	return &mod, repo.db.ConvertError(err)
}
