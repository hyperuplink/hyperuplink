package categories

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
)

func MoveUp(
	rt *runtime.Runtime,
	in *MoveInput,
) (err error) {
	cat := new(category.Category)
	if cat.ID, err = uuid.Parse(in.ID); err != nil {
		return err
	}

	return rt.Repositories.Category.MoveUp(cat)
}
