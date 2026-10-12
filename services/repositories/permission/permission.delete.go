package permission

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"xn--gckvb8fzb.com/glides/services/database"
)

func (repo *Repository) Remove(
	groupID pgtype.Text,
	categoryID pgtype.UUID,
) (err error) {
	_, err = database.Exec(context.Background(), repo.db, `DELETE FROM permissions
		WHERE group_id IS NOT DISTINCT FROM $1
		AND category_id IS NOT DISTINCT FROM $2`,
		groupID,
		categoryID,
	)
	return err
}
