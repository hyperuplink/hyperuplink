package forum

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mrusme/hyperuplink/models/vforum"
	"github.com/mrusme/hyperuplink/services/repositories/common"
)

func (repo *Repository) VAll(qo common.QueryOptions) (model *[]vforum.VForum, err error) {
	var rows pgx.Rows
	var mod []vforum.VForum

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM vforums`,
		common.QueryCapabilities{
			HasDeleted: true,
		},
	))
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[vforum.VForum])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) VAllForCategoryUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *[]vforum.VForum, err error) {
	var rows pgx.Rows
	var mod []vforum.VForum

	rows, err = repo.db.Query(
		qo.Query(`SELECT * FROM vforums
		WHERE category_id = $1`,
			common.QueryCapabilities{
				HasDeleted: true,
			},
		),
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[vforum.VForum])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) VAllForCategoryID(
	id string,
	qo common.QueryOptions,
) (model *[]vforum.VForum, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, repo.db.ConvertError(err)
	}

	return repo.VAllForCategoryUUID(uuID, qo)
}
