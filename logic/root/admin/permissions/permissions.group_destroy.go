package permissions

import (
	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/group"
)

func GroupDestroy(
	rt *runtime.Runtime,
	in *GroupDestroyInput,
) (err error) {
	grp := new(group.Group)
	grp.ID = in.ID

	return gh.Repositories(rt).Group.Delete(grp)
}
