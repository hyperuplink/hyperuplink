package forums

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/user"
	"github.com/mrusme/hyperuplink/models/vforum"
	"github.com/mrusme/hyperuplink/models/vtopic"
	"github.com/mrusme/hyperuplink/services/repositories/common"
)

func (r *Route) Show(c fiber.Ctx) (err error) {
	myRoute := route.For("CategoriesForums")
	req := request.New(r, c, myRoute,
		[]string{"base"}, "categories/forums/show",
		"")

	if ret, rerr := req.AccessControl(user.GuestRole); ret {
		return rerr
	}

	var fum *vforum.VForum
	fum, err = r.Runtime.Repositories.Forum.VGetBySlug(
		c.Params("forums"), // TODO: Abstract into req, automatic err handling
		common.QueryOptions{
			Limit: 1,
		},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.UpdateTitle(fum.Name)
	req.UpdateParentHref(req.HrefTo("_" + fum.CategorySlug))
	req.UpdateParentTitle(fum.CategoryName)
	req.SetData("forum", fum)

	var tops *[]vtopic.VTopic
	tops, err = r.Runtime.Repositories.Topic.VAllForForumUUID(
		fum.ID,
		common.QueryOptions{
			OrderBy: "updated_at",
			Order:   common.Descending,
		},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("topics", tops)

	return req.Respond()
}
