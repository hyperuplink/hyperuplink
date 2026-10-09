package permissions

import (
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/group"
)

func GroupDestroy(
	rt *runtime.Runtime,
	in *GroupDestroyInput,
) (err error) {
	grp := new(group.Group)
	grp.ID = in.ID

	if _, err = gh.Repositories(rt).Group.GetByID(grp.ID, common.QueryOptions{Limit: 1}); err != nil {
		return err
	}

	return gh.Repositories(rt).Group.Delete(grp)
}
