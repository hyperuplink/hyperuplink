package categories

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/category"
	"github.com/mrusme/hyperuplink/models/forum"
	"github.com/mrusme/hyperuplink/models/user"
	"github.com/mrusme/hyperuplink/services/repositories/common"
)

func (r *Route) Show(c fiber.Ctx) (err error) {
	myRoute := route.For("Categories")
	req := request.New(r, c, myRoute,
		[]string{"base"}, "categories/show",
		"")

	if ret, rerr := req.AccessControl(user.GuestRole); ret {
		return rerr
	}

	var cat *category.Category
	cat, err = r.Runtime.Repositories.Category.GetBySlug(
		c.Params("categories"), // TODO: Abstract into req, automatic err handling
		common.QueryOptions{
			Limit: 1,
		},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.UpdateTitle(cat.Name)
	req.SetData("category", cat)

	var fums *[]forum.Forum
	fums, err = r.Runtime.Repositories.Forum.AllForCategoryUUID(
		cat.ID,
		common.QueryOptions{
			OrderBy: "position",
			Order:   common.Ascending,
		},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("forums", fums)

	return req.Respond()
}
