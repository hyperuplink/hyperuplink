package topics

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/user"
	"github.com/mrusme/hyperuplink/models/vtopic"
	"github.com/mrusme/hyperuplink/services/repositories/common"
)

func (r *Route) Show(c fiber.Ctx) (err error) {
	myRoute := route.For("CategoriesForumsTopics")
	req := request.New(r, c, myRoute,
		[]string{"base"}, "categories/forums/topics/show",
		"")

	if ret, rerr := req.AccessControl(
		user.GuestRole, // TODO: Remove!
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	var top *vtopic.VTopic
	top, err = r.Runtime.Repositories.Topic.VGetBySlugs(
		c.Params("forums"), // TODO: Abstract into req, automatic err handling
		c.Params("topics"), // TODO: Abstract into req, automatic err handling
		common.QueryOptions{
			Limit: 1,
		},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.UpdateTitle(top.Name)
	req.UpdateParentHref(req.HrefTo(top.ForumSlug))
	req.UpdateParentTitle(top.ForumName)
	req.UpdateGrandParentHref(req.HrefTo("_" + top.CategorySlug))
	req.UpdateGrandParentTitle(top.CategoryName)

	top.Text, err = r.Runtime.Markdown.Convert(top.Text)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	req.SetData("topic", top)

	return req.Respond()
}
