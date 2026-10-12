package forums

import (
	"context"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
)

func Update(
	ctx context.Context,
	rt *runtime.Runtime,
	in *UpdateInput,
) (err error) {
	fum := new(forum.Forum)
	if fum.ID, err = uuid.Parse(in.ID); err != nil {
		return err
	}
	fum.Name = in.Name
	fum.Slug = in.Slug
	if fum.CategoryID, err = uuid.Parse(in.CategoryID); err != nil {
		return err
	}
	fum.Description = in.Description

	if _, err = gh.Repositories(rt).Forum.GetByUUID(ctx, fum.ID, common.QueryOptions{Limit: 1}); err != nil {
		return err
	}

	return gh.Repositories(rt).Forum.Update(ctx, fum)
}
