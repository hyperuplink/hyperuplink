package contact

import (
	"errors"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicdocs "xn--gckvb8fzb.com/hyperuplink/logic/root/docs"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("DocsContact")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.GuestRole,
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	var view *logicdocs.Page
	view, err = logicdocs.Show(r.Runtime, logicdocs.PageContact)
	if errors.Is(err, errs.ErrNoRows) {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("html", view.HTML)

	return req.Respond()
}
