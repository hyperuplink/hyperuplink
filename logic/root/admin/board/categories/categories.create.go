package categories

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
)

func Create(
	rt *runtime.Runtime,
	in *CreateInput,
) (id uuid.UUID, err error) {
	cat := new(category.Category)
	cat.Name = in.Name
	cat.Slug = in.Slug

	return rt.Repositories.Category.Create(cat)
}
