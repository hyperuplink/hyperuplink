package permission

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
)

func (repo *Repository) All() (model []permission.Permission, err error) {
	return database.Rows[permission.Permission](context.Background(), repo.db, `SELECT * FROM permissions
		WHERE deleted_at IS NULL`)
}

func (repo *Repository) AllForGroup(
	groupID string,
) (model []permission.Permission, err error) {
	return database.Rows[permission.Permission](context.Background(), repo.db, `SELECT * FROM permissions
		WHERE group_id = $1
		AND deleted_at IS NULL`,
		groupID,
	)
}

func (repo *Repository) GetDefault() (model *permission.Permission, err error) {
	return repo.GetFor(pgtype.Text{}, pgtype.UUID{})
}

func (repo *Repository) GetFor(
	groupID pgtype.Text,
	categoryID pgtype.UUID,
) (model *permission.Permission, err error) {
	return database.Row[permission.Permission](context.Background(), repo.db, `SELECT * FROM permissions
		WHERE group_id IS NOT DISTINCT FROM $1
		AND category_id IS NOT DISTINCT FROM $2
		AND deleted_at IS NULL
		LIMIT 1`,
		groupID,
		categoryID,
	)
}
