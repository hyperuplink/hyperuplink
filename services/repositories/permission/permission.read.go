package permission

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mrusme/hyperuplink/models/permission"
	"github.com/mrusme/hyperuplink/models/user"
)

func (repo *Repository) All() (model *[]permission.Permission, err error) {
	var rows pgx.Rows
	var mod []permission.Permission

	rows, err = repo.db.Query(`SELECT * FROM permissions
		WHERE deleted_at IS NULL
		`)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[permission.Permission])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetByUUID(
	id uuid.UUID,
) (model *permission.Permission, err error) {
	var rows pgx.Rows
	var mod permission.Permission

	rows, err = repo.db.Query(`SELECT * FROM permissions
		WHERE id = $1
		AND deleted_at IS NULL
		LIMIT 1`,
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[permission.Permission])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetByID(
	id string,
) (model *permission.Permission, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, repo.db.ConvertError(err)
	}

	return repo.GetByUUID(uuID)
}

func (repo *Repository) GetByRoleUnitForumUUID(
	role user.Role,
	unit pgtype.Text,
	id uuid.UUID,
) (model *permission.Permission, err error) {
	var rows pgx.Rows
	var mod permission.Permission

	rows, err = repo.db.Query(`SELECT * FROM permissions
		WHERE role = $1
		AND unit = $2
		AND forum_id = $3
		AND deleted_at IS NULL
		LIMIT 1`,
		role,
		unit,
		id,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[permission.Permission])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetByRoleUnitForumID(
	role user.Role,
	unit pgtype.Text,
	id string,
) (model *permission.Permission, err error) {
	var uuID uuid.UUID

	if uuID, err = uuid.Parse(id); err != nil {
		return nil, repo.db.ConvertError(err)
	}

	return repo.GetByRoleUnitForumUUID(role, unit, uuID)
}
