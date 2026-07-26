package forums

import (
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
)

func View(rt *runtime.Runtime) (catsfums []CategoryWithForums, err error) {
	cats, err := gh.Repositories(rt).Category.All(common.QueryOptions{
		OrderBy: "position",
		Order:   common.Ascending,
	})
	if err != nil {
		return nil, err
	}

	fums, err := gh.Repositories(rt).Forum.All(common.QueryOptions{
		OrderBy: "position",
		Order:   common.Ascending,
	})
	if err != nil {
		return nil, err
	}

	for _, cat := range *cats {
		catfum := CategoryWithForums{
			Category: cat,
		}
		for _, fum := range *fums {
			if fum.CategoryID == cat.ID {
				catfum.Forums = append(catfum.Forums, fum)
			}
		}
		catsfums = append(catsfums, catfum)
	}

	return catsfums, nil
}
