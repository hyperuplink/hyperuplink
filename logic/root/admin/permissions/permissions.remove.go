package permissions

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
)

func Remove(
	ctx context.Context,
	rt *runtime.Runtime,
	in *RemoveInput,
) (err error) {
	uid, perr := uuid.Parse(in.CategoryID)
	if perr != nil {
		return perr
	}

	var groupID pgtype.Text
	if in.GroupID != "" {
		groupID = pgtype.Text{String: in.GroupID, Valid: true}
	}
	categoryID := pgtype.UUID{Bytes: [16]byte(uid), Valid: true}

	return gh.Repositories(rt).Permission.Remove(ctx, groupID, categoryID)
}
