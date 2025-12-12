package topics

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/helpers"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/http/web/request/site"
	"github.com/mrusme/hyperuplink/models/user"
	"github.com/mrusme/hyperuplink/models/vreply"
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
	req.UpdateParentHref(req.HrefTo(route.For("CategoriesForums").Fill(
		map[string]string{
			"categories": top.CategorySlug,
			"forums":     top.ForumSlug,
		},
	).AsURL()))
	req.UpdateParentTitle(top.ForumName)
	req.UpdateGrandParentHref(req.HrefTo(route.For("Categories").Fill(
		map[string]string{
			"categories": top.CategorySlug,
		},
	).AsURL()))
	req.UpdateGrandParentTitle(top.CategoryName)

	// TODO: Move to Create/Update and save Markdown
	top.Text, err = r.Runtime.Markdown.Convert(top.Text)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	req.SetData("topic", top)

	page, err := strconv.Atoi(c.Query("page", "1"))
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var reps *[]vreply.VReply
	var total int64
	var perPage int = 2
	reps, total, err = r.Runtime.Repositories.Reply.VAllForTopicUUID(
		top.ID,
		common.QueryOptions{
			OrderBy: "created_at",
			Order:   common.Ascending,
			Limit:   perPage,
			Page:    page,
		},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	pages := helpers.GetNumberOfPages(total, perPage)

	req.Site.SetPager(site.NewPager(pages, perPage, page))

	req.SetData("replies", reps)

	return req.Respond()
}
