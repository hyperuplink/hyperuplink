package categories

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
)

func MoveDown(
	rt *runtime.Runtime,
	in *MoveInput,
) (err error) {
	cat := new(category.Category)
	if cat.ID, err = uuid.Parse(in.ID); err != nil {
		return err
	}

	if _, err = gh.Repositories(rt).Category.GetByUUID(cat.ID, common.QueryOptions{Limit: 1}); err != nil {
		return err
	}

	return gh.Repositories(rt).Category.MoveDown(cat)
}
