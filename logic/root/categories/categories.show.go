package categories

import (
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/vforum"
)

type View struct {
	Category *category.Category `json:"category"`
	Forums   *[]vforum.VForum   `json:"forums"`
}

func Show(
	rt *runtime.Runtime,
	slug string,
	perms *permission.Resolution,
) (view *View, err error) {
	cat, err := gh.Repositories(rt).Category.GetBySlug(
		slug,
		common.QueryOptions{
			Limit: 1,
		},
	)
	if err != nil {
		return nil, err
	}

	if !perms.CanReadID(cat.ID) {
		return nil, errs.ErrForbidden
	}

	fums, err := gh.Repositories(rt).Forum.VAllForCategoryUUID(
		cat.ID,
		common.QueryOptions{
			OrderBy: "position",
			Order:   common.Ascending,
		},
	)
	if err != nil {
		return nil, err
	}

	return &View{
		Category: cat,
		Forums:   fums,
	}, nil
}
