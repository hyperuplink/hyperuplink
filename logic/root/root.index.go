package root

import (
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/vforum"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
)

type CategoryWithForums struct {
	Category category.Category `json:"category"`
	Forums   []vforum.VForum   `json:"forums"`
}

type Board struct {
	CategoriesForums []CategoryWithForums `json:"categories_forums"`
	RecentTopics     []vtopic.VTopic      `json:"recent_topics"`
}

func BoardView(
	rt *runtime.Runtime,
	perms *permission.Resolution,
) (board *Board, err error) {
	cats, err := gh.Repositories(rt).Category.All(common.QueryOptions{
		OrderBy: "position",
		Order:   common.Ascending,
	})
	if err != nil {
		return nil, err
	}

	fums, err := gh.Repositories(rt).Forum.VAll(common.QueryOptions{
		OrderBy: "position",
		Order:   common.Ascending,
	})
	if err != nil {
		return nil, err
	}

	board = new(Board)
	for _, cat := range *cats {
		if !perms.CanReadID(cat.ID) {
			continue
		}
		catfum := CategoryWithForums{
			Category: cat,
		}
		for _, fum := range *fums {
			if fum.CategoryID == cat.ID {
				catfum.Forums = append(catfum.Forums, fum)
			}
		}
		board.CategoriesForums = append(board.CategoriesForums, catfum)
	}

	tops, err := gh.Repositories(rt).Topic.VAll(
		common.QueryOptions{
			OrderBy: "updated_at",
			Order:   common.Descending,
			Limit:   5,
		},
	)
	if err != nil {
		return nil, err
	}

	board.RecentTopics = []vtopic.VTopic{}
	for _, top := range *tops {
		if perms.CanReadSlug(top.CategorySlug) {
			board.RecentTopics = append(board.RecentTopics, top)
		}
	}

	return board, nil
}
