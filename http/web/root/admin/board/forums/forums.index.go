package forums

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

type CategoryWithForums struct {
	Category category.Category
	Forums   []forum.Forum
}

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardForums")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	var cats *[]category.Category
	cats, err = r.Runtime.Repositories.Category.All(common.QueryOptions{
		OrderBy: "position",
		Order:   common.Ascending,
	})
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var fums *[]forum.Forum
	fums, err = r.Runtime.Repositories.Forum.All(common.QueryOptions{
		OrderBy: "position",
		Order:   common.Ascending,
	})
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var catsfums []CategoryWithForums
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

	req.SetData("categories_forums", catsfums)

	return req.Respond()
}
