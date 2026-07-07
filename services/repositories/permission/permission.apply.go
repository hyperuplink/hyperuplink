package permission

import (
	"github.com/jackc/pgx/v5/pgtype"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
)

func (repo *Repository) Apply(
	groupID pgtype.Text,
	categoryID pgtype.UUID,
	level byte,
) (err error) {
	bits := permission.LevelToBitString(level)

	tx, err := repo.db.Tx()
	if err != nil {
		return repo.db.ConvertError(err)
	}
	defer tx.End()

	tag, err := tx.Exec(`UPDATE permissions SET
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
		return repo.db.ConvertError(err)
	}

	if tag.RowsAffected() == 0 {
		_, err = tx.Exec(`INSERT INTO permissions (
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
			return repo.db.ConvertError(err)
		}
	}

	return repo.db.ConvertError(tx.Commit())
}
