package permission

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"xn--gckvb8fzb.com/glides/services/database"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (repo *Repository) Resolve(
	role user.Role,
	memberOf []string,
	cats []category.Category,
) (res *permission.Resolution, err error) {
	res = permission.NewResolution()

	if role == user.AdminRole {
		res.IsAdmin = true
		return res, nil
	}

	perms, err := database.Rows[permission.Permission](context.Background(), repo.db, `SELECT * FROM permissions
		WHERE (group_id IS NULL OR group_id = ANY($1))
		AND deleted_at IS NULL`,
		memberOf,
	)
	if err != nil {
		return nil, err
	}

	categoryDefault := map[uuid.UUID]byte{}
	groupMax := map[uuid.UUID]byte{}
	for i := range perms {
		p := &perms[i]
		lvl := p.Level()
		switch {
		case !p.GroupID.Valid && !p.CategoryID.Valid:
			res.Default = lvl
		case !p.GroupID.Valid:
			categoryDefault[uuid.UUID(p.CategoryID.Bytes)] = lvl
		case p.CategoryID.Valid:
			id := uuid.UUID(p.CategoryID.Bytes)
			if lvl > groupMax[id] {
				groupMax[id] = lvl
			}
		}
	}

	if cats != nil {
		for _, cat := range cats {
			level := res.Default
			if lvl, ok := categoryDefault[cat.ID]; ok {
				level = lvl
			}
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

	catDef, err := repo.GetFor(pgtype.Text{}, pgtype.UUID{Bytes: categoryID, Valid: true})
	if err != nil && !errors.Is(err, errs.ErrNoRows) {
		return permission.None, err
	}
	if err == nil && catDef != nil {
		level = catDef.Level()
	}

	if len(memberOf) == 0 {
		return level, nil
	}

	perms, err := database.Rows[permission.Permission](context.Background(), repo.db, `SELECT * FROM permissions
		WHERE group_id = ANY($1)
		AND category_id = $2
		AND deleted_at IS NULL`,
		memberOf,
		categoryID,
	)
	if err != nil {
		return level, err
	}

	for i := range perms {
		if lvl := perms[i].Level(); lvl > level {
			level = lvl
		}
	}

	return level, nil
}
