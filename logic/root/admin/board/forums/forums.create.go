package forums

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
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

	return rt.Repositories.Forum.Create(fum)
}
