package categories

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
)

func Update(
	rt *runtime.Runtime,
	in *UpdateInput,
) (err error) {
	cat := new(category.Category)
	if cat.ID, err = uuid.Parse(in.ID); err != nil {
		return err
	}
	cat.Name = in.Name
	cat.Slug = in.Slug

	return rt.Repositories.Category.Update(cat)
}
