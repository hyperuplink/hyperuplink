package forums

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request/site"
	logicforums "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
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

	var activePage int
	activePage, err = strconv.Atoi(c.Query("page", "1"))
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	viewerID := uuid.NullUUID{}
	viewerID.UUID, viewerID.Valid = req.Session.GetUserUUID()

	var perPage int = req.System.GetTopicsPerPage()

	var view *logicforums.View
	view, err = logicforums.Show(r.Runtime, &logicforums.ShowInput{
		ForumSlug: c.Params("forums"),
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
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	fum := view.Forum

	req.Menu.SetCategoryForumSlugs(fum.CategorySlug, fum.Slug)

	req.UpdateTitle(fum.Name)
	req.UpdateParentHref(req.HrefTo("_" + fum.CategorySlug))
	req.UpdateParentTitle(fum.CategoryName)
	req.SetData("forum", fum)

	req.Site.SetPager(site.NewPager(view.Pages, perPage, activePage))

	req.SetData("topics", view.Topics)

	if view.Unread != nil {
		req.SetData("unread", view.Unread)
	}

	return req.Respond()
}
