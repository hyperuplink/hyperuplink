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

	var reps *[]vreply.VReply
	var activePage int
	var total int64
	var perPage int = 2 // TODO: Move to System
	var limit int = perPage
	var offAdjust int = 1

	activePage, err = strconv.Atoi(c.Query("page", "1"))
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	// If we are on the first page, we subtract 1 from limit due to the Topic,
	// and we set offadjust to 0 because we only need to adjust the offset for
	// all pages past the first one.
	if activePage == 1 {
		limit -= 1
		offAdjust = 0
	}
	reps, total, err = r.Runtime.Repositories.Reply.VAllForTopicUUID(
		top.ID,
		common.QueryOptions{
			OrderBy:   "created_at",
			Order:     common.Ascending,
			Limit:     limit,
			Page:      activePage,
			OffAdjust: offAdjust,
		},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	// We add the Topic to the total
	total += 1
	pages := helpers.GetNumberOfPages(total, perPage)
	req.Site.SetPager(site.NewPager(pages, perPage, activePage))

	req.SetData("replies", reps)

	replyTo := c.Query("reply", "")
	if replyTo != "" {
		for _, rep := range *reps {
			if rep.ID.String() == replyTo {
				req.SetData("reply_to", rep)
			}
		}
	}

	return req.Respond()
}
