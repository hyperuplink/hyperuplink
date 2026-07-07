package permission

import (
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
)

func (repo *Repository) All() (model *[]permission.Permission, err error) {
	var rows pgx.Rows
	var mod []permission.Permission

	rows, err = repo.db.Query(`SELECT * FROM permissions
		WHERE deleted_at IS NULL`)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[permission.Permission])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) AllForGroup(
	groupID string,
) (model *[]permission.Permission, err error) {
	var rows pgx.Rows
	var mod []permission.Permission

	rows, err = repo.db.Query(`SELECT * FROM permissions
		WHERE group_id = $1
		AND deleted_at IS NULL`,
		groupID,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectRows(rows, pgx.RowToStructByName[permission.Permission])

	return &mod, repo.db.ConvertError(err)
}

func (repo *Repository) GetDefault() (model *permission.Permission, err error) {
	return repo.GetFor(pgtype.Text{}, pgtype.UUID{})
}

func (repo *Repository) GetFor(
	groupID pgtype.Text,
	categoryID pgtype.UUID,
) (model *permission.Permission, err error) {
	var rows pgx.Rows
	var mod permission.Permission

	rows, err = repo.db.Query(`SELECT * FROM permissions
		WHERE group_id IS NOT DISTINCT FROM $1
		AND category_id IS NOT DISTINCT FROM $2
		AND deleted_at IS NULL
		LIMIT 1`,
		groupID,
		categoryID,
	)
	if err != nil {
		return nil, repo.db.ConvertError(err)
	}

	mod, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[permission.Permission])

	return &mod, repo.db.ConvertError(err)
}
