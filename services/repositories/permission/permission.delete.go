package permission

import (
	"github.com/jackc/pgx/v5/pgtype"
)

func (repo *Repository) Remove(
	groupID pgtype.Text,
	categoryID pgtype.UUID,
) (err error) {
	_, err = repo.db.Exec(`DELETE FROM permissions
		WHERE group_id IS NOT DISTINCT FROM $1
		AND category_id IS NOT DISTINCT FROM $2`,
		groupID,
		categoryID,
	)
	return repo.db.ConvertError(err)
}
