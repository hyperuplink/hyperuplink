package permissions

import (
	"xn--gckvb8fzb.com/hyperuplink/models/group"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
)

func GroupDestroy(
	rt *runtime.Runtime,
	in *GroupDestroyInput,
) (err error) {
	grp := new(group.Group)
	grp.ID = in.ID

	return rt.Repositories.Group.Delete(grp)
}
