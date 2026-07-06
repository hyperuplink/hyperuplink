package unit

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"xn--gckvb8fzb.com/hyperuplink/models/unit"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

func (repo *Repository) GetByUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *unit.Unit, err error) {
	var rows pgx.Rows
	var mod unit.Unit

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM units WHERE id = $1`,
		common.QueryCapabilities{
			HasBanned:  true,
			HasDeleted: true,
		}),
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[unit.Unit])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetByID(
	id string,
	qo common.QueryOptions,
) (model *unit.Unit, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, repo.db.ConvertError(err)
	}

	return repo.GetByUUID(uuID, qo)
}
