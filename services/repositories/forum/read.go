package forum

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/mrusme/hyperuplink/models/forum"
)

func (repo *Repository) AllForCategoryUUID(
	id uuid.UUID,
) (model *[]forum.Forum, err error) {
	var rows pgx.Rows
	var mod []forum.Forum

	rows, err = repo.db.Query(`SELECT * FROM forums
		WHERE category_id = $1
		AND deleted_at IS NULL
		LIMIT 1`,
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[forum.Forum])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) AllForCategoryID(
	id string,
) (model *[]forum.Forum, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, repo.db.ConvertError(err)
	}

	return repo.AllForCategoryUUID(uuID)
}

func (repo *Repository) GetByUUID(
	id uuid.UUID,
) (model *forum.Forum, err error) {
	var rows pgx.Rows
	var mod forum.Forum

	rows, err = repo.db.Query(`SELECT * FROM forums
		WHERE id = $1
		AND deleted_at IS NULL
		LIMIT 1`,
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[forum.Forum])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetByID(
	id string,
) (model *forum.Forum, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, repo.db.ConvertError(err)
	}

	return repo.GetByUUID(uuID)
}
