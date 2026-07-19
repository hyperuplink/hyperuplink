package categories

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
)

func Destroy(
	rt *runtime.Runtime,
	in *DestroyInput,
) (err error) {
	cat := new(category.Category)
	if cat.ID, err = uuid.Parse(in.ID); err != nil {
		return err
	}

	return rt.Repositories.Category.Delete(cat)
}
