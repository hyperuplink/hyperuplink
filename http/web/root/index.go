package root

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/category"
	"github.com/mrusme/hyperuplink/models/user"
	"github.com/mrusme/hyperuplink/models/vforum"
	"github.com/mrusme/hyperuplink/models/vtopic"
	"github.com/mrusme/hyperuplink/services/repositories/common"
)

type CategoryWithForums struct {
	Category category.Category
	Forums   []vforum.VForum
}

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("Root")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL(),
		"")

	if ret, rerr := req.AccessControl(
		user.GuestRole,
		user.UserRole,
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

	var fums *[]vforum.VForum
	fums, err = r.Runtime.Repositories.Forum.VAll(common.QueryOptions{
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

	var tops *[]vtopic.VTopic
	tops, err = r.Runtime.Repositories.Topic.VAll(
		common.QueryOptions{
			OrderBy: "updated_at",
			Order:   common.Descending,
			Limit:   5,
		},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("topics", tops)

	return req.Respond()
}
