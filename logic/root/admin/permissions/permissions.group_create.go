package permissions

import (
	"xn--gckvb8fzb.com/hyperuplink/models/group"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
)

func GroupCreate(
	rt *runtime.Runtime,
	in *GroupCreateInput,
) (err error) {
	grp := new(group.Group)
	grp.ID = in.ID
	grp.Name = in.Name

	_, err = rt.Repositories.Group.Create(grp)

	return err
}
