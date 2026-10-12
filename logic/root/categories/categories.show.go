package categories

import (
	"context"
	"github.com/google/uuid"
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
	Forums   []vforum.VForum    `json:"forums"`
}

func Resolve(
	ctx context.Context,
	rt *runtime.Runtime,
	idOrSlug string,
) (cat *category.Category, err error) {
	if id, perr := uuid.Parse(idOrSlug); perr == nil {
		return gh.Repositories(rt).Category.GetByUUID(ctx, id, common.QueryOptions{Limit: 1})
	}

	return gh.Repositories(rt).Category.GetBySlug(ctx, idOrSlug, common.QueryOptions{Limit: 1})
}

func Show(
	ctx context.Context,
	rt *runtime.Runtime,
	idOrSlug string,
	perms *permission.Resolution,
) (view *View, err error) {
	cat, err := Resolve(ctx, rt, idOrSlug)
	if err != nil {
		return nil, err
	}

	if !perms.CanReadID(cat.ID) {
		return nil, errs.ErrForbidden
	}

	fums, err := gh.Repositories(rt).Forum.VAllForCategoryUUID(ctx,
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
