package forums

import (
	"context"
	"slices"

	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/vforum"
)

func Index(
	ctx context.Context,
	rt *runtime.Runtime,
	perms *permission.Resolution,
) (forums []vforum.VForum, err error) {
	cats, err := gh.Repositories(rt).Category.All(ctx, common.QueryOptions{
		OrderBy: "position",
		Order:   common.Ascending,
	})
	if err != nil {
		return nil, err
	}

	fums, err := gh.Repositories(rt).Forum.VAll(ctx, common.QueryOptions{
		OrderBy: "position",
		Order:   common.Ascending,
	})
	if err != nil {
		return nil, err
	}

	positions := make(map[uuid.UUID]int, len(cats))
	for _, cat := range cats {
		positions[cat.ID] = cat.Position
	}

	forums = make([]vforum.VForum, 0, len(fums))
	for _, fum := range fums {
		if perms.CanReadID(fum.CategoryID) {
			forums = append(forums, fum)
		}
	}

	slices.SortStableFunc(forums, func(a, b vforum.VForum) int {
		if positions[a.CategoryID] != positions[b.CategoryID] {
			return positions[a.CategoryID] - positions[b.CategoryID]
		}
		return a.Position - b.Position
	})

	return forums, nil
}
