package categories

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/models/vforum"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

func (r *Route) Show(c fiber.Ctx) (err error) {
	myRoute := route.For("Categories")
	req := request.New(r, c, myRoute,
		[]string{"base"}, "categories/show",
		"")

	if ret, rerr := req.AccessControl(
		user.GuestRole, // TODO: Remove!
		user.UserRole,
		user.AdminRole,
	); ret {
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

	var fums *[]vforum.VForum
	fums, err = r.Runtime.Repositories.Forum.VAllForCategoryUUID(
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
