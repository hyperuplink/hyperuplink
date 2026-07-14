package forums

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request/site"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/models/vforum"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

func (r *Route) Show(c fiber.Ctx) (err error) {
	myRoute := route.For("CategoriesForums")
	req := request.New(r, c, myRoute,
		[]string{"base"}, "categories/forums/show",
		"")

	if ret, rerr := req.AccessControl(
		user.GuestRole, // TODO: Remove!
		user.UserRole,
		user.AdminRole,
	); ret {
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

	if !req.Perms().CanReadID(fum.CategoryID) {
		return req.RedirectToRoot()
	}

	req.Menu.SetCategoryForumSlugs(fum.CategorySlug, fum.Slug)

	req.UpdateTitle(fum.Name)
	req.UpdateParentHref(req.HrefTo("_" + fum.CategorySlug))
	req.UpdateParentTitle(fum.CategoryName)
	req.SetData("forum", fum)

	var activePage int
	var total int64
	var perPage int = req.System.GetTopicsPerPage()
	var limit int = perPage

	activePage, err = strconv.Atoi(c.Query("page", "1"))
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var tops *[]vtopic.VTopic
	tops, total, err = r.Runtime.Repositories.Topic.VAllForForumUUID(
		fum.ID,
		common.QueryOptions{
			OrderBy: "updated_at",
			Order:   common.Descending,
			Limit:   limit,
			Page:    activePage,
		},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	pages := helpers.GetNumberOfPages(total, perPage)
	req.Site.SetPager(site.NewPager(pages, perPage, activePage))

	req.SetData("topics", tops)

	if actorID, ok := req.Session.GetUserUUID(); ok {
		var unread map[string]bool
		unread, err = logicactivity.UnreadTopics(r.Runtime, actorID, tops)
		if ret, rerr := req.RespondOnError(err); ret == true {
			return rerr
		}
		req.SetData("unread", unread)
	}

	return req.Respond()
}
