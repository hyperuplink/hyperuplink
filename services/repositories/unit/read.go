package unit

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mrusme/hyperuplink/models/unit"
)

func (repo *Repository) GetByUUID(
	id uuid.UUID,
) (model *unit.Unit, err error) {
	var rows pgx.Rows
	var mod unit.Unit

	rows, err = repo.db.Query(`SELECT * FROM units
		WHERE id = $1
		AND banned_at IS NULL
		AND deleted_at IS NULL
		LIMIT 1`,
		id)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[unit.Unit])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetByID(
	id string,
) (model *unit.Unit, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, repo.db.ConvertError(err)
	}

	return repo.GetByUUID(uuID)
}
