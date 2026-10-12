package permission

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
)

func (repo *Repository) Apply(
	ctx context.Context,
	groupID pgtype.Text,
	categoryID pgtype.UUID,
	level byte,
) (err error) {
	bits := permission.LevelToBitString(level)

	tx, err := repo.db.Tx(ctx)
	if err != nil {
		return database.ConvertError(err)
	}
	defer tx.End(ctx)

	tag, err := database.Exec(ctx, tx, `UPDATE permissions SET
		 bits =       $3::BIT(3)
		,updated_at = NOW()
		WHERE group_id IS NOT DISTINCT FROM $1
		AND category_id IS NOT DISTINCT FROM $2
		AND deleted_at IS NULL`,
		groupID,
		categoryID,
		bits,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		_, err = database.Exec(ctx, tx, `INSERT INTO permissions (
			 group_id
			,category_id
			,bits
			) VALUES (
			 $1
			,$2
			,$3::BIT(3)
			)`,
			groupID,
			categoryID,
			bits,
		)
		if err != nil {
			return err
		}
	}

	return database.ConvertError(tx.Commit(ctx))
}
