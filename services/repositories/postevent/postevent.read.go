package postevent

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mrusme/hyperuplink/models/postevent"
	"github.com/mrusme/hyperuplink/services/repositories/common"
)

func (repo *Repository) GetByUUID(
	id uuid.UUID,
	qo common.QueryOptions,
) (model *postevent.PostEvent, err error) {
	var rows pgx.Rows
	var mod postevent.PostEvent

	rows, err = repo.db.Query(qo.Query(
		`SELECT * FROM postevents WHERE id = $1`,
		common.QueryCapabilities{
			HasDeleted: true,
		}),
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[postevent.PostEvent])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetByID(
	id string,
	qo common.QueryOptions,
) (model *postevent.PostEvent, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, repo.db.ConvertError(err)
	}

	return repo.GetByUUID(uuID, qo)
}
