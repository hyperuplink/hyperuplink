package forums

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
)

func Update(
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

	return gh.Repositories(rt).Forum.Update(fum)
}
