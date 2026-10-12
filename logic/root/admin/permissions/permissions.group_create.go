package permissions

import (
	"context"
	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/group"
)

func GroupCreate(
	ctx context.Context,
	rt *runtime.Runtime,
	in *GroupCreateInput,
) (err error) {
	grp := new(group.Group)
	grp.ID = in.ID
	grp.Name = in.Name

	_, err = gh.Repositories(rt).Group.Create(ctx, grp)

	return err
}
