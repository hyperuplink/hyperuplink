package categories

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logiccategories "xn--gckvb8fzb.com/hyperuplink/logic/root/categories"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
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

	var view *logiccategories.View
	view, err = logiccategories.Show(r.Runtime, c.Params("categories"), req.Perms())
	if errors.Is(err, errs.ErrNoRows) {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if errors.Is(err, errs.ErrForbidden) {
		return req.RedirectToRoot()
	}
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.UpdateTitle(view.Category.Name)
	req.SetData("category", view.Category)
	req.SetData("forums", view.Forums)

	return req.Respond()
}
