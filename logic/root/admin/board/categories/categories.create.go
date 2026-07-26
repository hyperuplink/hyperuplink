package categories

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
)

func Create(
	rt *runtime.Runtime,
	in *CreateInput,
) (id uuid.UUID, err error) {
	cat := new(category.Category)
	cat.Name = in.Name
	cat.Slug = in.Slug

	return gh.Repositories(rt).Category.Create(cat)
}
