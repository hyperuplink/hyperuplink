package permission

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
)

func (repo *Repository) EffectiveLevel(
	memberOf []string,
	categoryID uuid.UUID,
) (level byte, err error) {
	level = permission.None

	def, err := repo.GetDefault()
	if err != nil && !errors.Is(err, errs.ErrNoRows) {
		return permission.None, err
	}
	if err == nil && def != nil {
		level = def.Level()
	}

	if len(memberOf) == 0 {
		return level, nil
	}

	var rows pgx.Rows
	rows, err = repo.db.Query(`SELECT * FROM permissions
		WHERE group_id = ANY($1)
		AND category_id = $2
		AND deleted_at IS NULL`,
		memberOf,
		categoryID,
	)
	if err != nil {
		return level, repo.db.ConvertError(err)
	}

	perms, err := pgx.CollectRows(rows, pgx.RowToStructByName[permission.Permission])
	if err != nil {
		return level, repo.db.ConvertError(err)
	}

	for i := range perms {
		if lvl := perms[i].Level(); lvl > level {
			level = lvl
		}
	}

	return level, nil
}
