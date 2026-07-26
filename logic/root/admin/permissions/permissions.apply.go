package permissions

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
)

func Apply(
	rt *runtime.Runtime,
	in *ApplyInput,
) (err error) {
	isDefault := in.GroupID == "" && in.CategoryID == ""
	isGroupMapping := in.GroupID != "" && in.CategoryID != ""

	if isDefault == false && isGroupMapping == false {
		return errs.ErrValidation
	}

	if isGroupMapping && in.Level == permission.LevelNone {
		return errs.ErrValidation
	}

	level, ok := permission.LevelFromString(in.Level)
	if ok == false {
		return errs.ErrValidation
	}

	var groupID pgtype.Text
	if in.GroupID != "" {
		groupID = pgtype.Text{String: in.GroupID, Valid: true}
	}

	var categoryID pgtype.UUID
	if in.CategoryID != "" {
		uid, perr := uuid.Parse(in.CategoryID)
		if perr != nil {
			return perr
		}
		categoryID = pgtype.UUID{Bytes: [16]byte(uid), Valid: true}
	}

	return gh.Repositories(rt).Permission.Apply(groupID, categoryID, level)
}
