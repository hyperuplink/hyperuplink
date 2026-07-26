package topics

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request/site"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums/topics"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
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

	var activePage int
	activePage, err = strconv.Atoi(c.Query("page", "1"))
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	viewerID := uuid.NullUUID{}
	viewerID.UUID, viewerID.Valid = req.Session.GetUserUUID()

	var perPage int = req.System.GetPostsPerPage()

	var view *logictopics.View
	view, err = logictopics.Show(r.Runtime, &logictopics.ShowInput{
		ForumSlug: c.Params("forums"),
		TopicSlug: c.Params("topics"),
		Page:      activePage,
		PerPage:   perPage,
		ViewerID:  viewerID,
	}, req.Perms())
	if errors.Is(err, errs.ErrNoRows) {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if errors.Is(err, errs.ErrForbidden) {
		return req.RedirectToRoot()
	}
	if ret, rerr := req.RedirectToRootOnError(err); ret == true {
		return rerr
	}

	top := view.Topic

	req.Menu.SetCategoryForumSlugs(top.CategorySlug, top.ForumSlug)

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

	req.SetData("topic", top)
	req.SetData("poll", view.Poll)

	req.Site.SetPager(site.NewPager(view.Pages, perPage, activePage))

	reps := view.Replies
	req.SetData("replies", reps)

	replyTo := c.Query("reply", "")
	if replyTo != "" {
		if top.ShortID == replyTo {
			req.SetData("reply_to", top)
			req.SetData("reply_to_type", "topic")
		} else {
			for _, rep := range *reps {
				if rep.ShortID == replyTo {
					req.SetData("reply_to", rep)
					req.SetData("reply_to_type", "reply")
				}
			}
		}
	}

	return req.Respond()
}
