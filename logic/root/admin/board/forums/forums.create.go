package forums

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
)

func Create(
	rt *runtime.Runtime,
	in *CreateInput,
) (id uuid.UUID, err error) {
	fum := new(forum.Forum)
	fum.Name = in.Name
	fum.Slug = in.Slug
	if fum.CategoryID, err = uuid.Parse(in.CategoryID); err != nil {
		return uuid.Nil, err
	}

	return gh.Repositories(rt).Forum.Create(fum)
}
