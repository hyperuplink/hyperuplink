package categories

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
)

func MoveUp(
	rt *runtime.Runtime,
	in *MoveInput,
) (err error) {
	cat := new(category.Category)
	if cat.ID, err = uuid.Parse(in.ID); err != nil {
		return err
	}

	return gh.Repositories(rt).Category.MoveUp(cat)
}
