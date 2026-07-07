package permission

import (
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (repo *Repository) Resolve(
	role user.Role,
	memberOf []string,
	cats *[]category.Category,
) (res *permission.Resolution, err error) {
	res = permission.NewResolution()

	if role == user.AdminRole {
		res.IsAdmin = true
		return res, nil
	}

	def, err := repo.GetDefault()
	if err != nil && !errors.Is(err, errs.ErrNoRows) {
		return nil, err
	}
	if err == nil && def != nil {
		res.Default = def.Level()
	}

	groupMax := map[uuid.UUID]byte{}
	if len(memberOf) > 0 {
		var rows pgx.Rows
		rows, err = repo.db.Query(`SELECT * FROM permissions
			WHERE group_id = ANY($1)
			AND category_id IS NOT NULL
			AND deleted_at IS NULL`,
			memberOf,
		)
		if err != nil {
			return nil, repo.db.ConvertError(err)
		}

		perms, cerr := pgx.CollectRows(rows, pgx.RowToStructByName[permission.Permission])
		if cerr != nil {
			return nil, repo.db.ConvertError(cerr)
		}

		for i := range perms {
			id := uuid.UUID(perms[i].CategoryID.Bytes)
			if lvl := perms[i].Level(); lvl > groupMax[id] {
				groupMax[id] = lvl
			}
		}
	}

	if cats != nil {
		for _, cat := range *cats {
			level := res.Default
			if lvl, ok := groupMax[cat.ID]; ok && lvl > level {
				level = lvl
			}
			res.SetCategory(cat.ID, cat.Slug, level)
		}
	}

	return res, nil
}

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
