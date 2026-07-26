package session

import (
	"github.com/gofiber/fiber/v3"
	goth_fiber "github.com/shareed2k/goth_fiber/v2"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) SignOutShow(c fiber.Ctx) (err error) {
	myRoute := route.For("SessionSignout")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL(),
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	provider, ok := req.Session.GetProvider()
	if ok == false {
		return req.RedirectToRoot()
	}

	if provider == "local" {
		if err = req.Session.Reset(); err != nil {
			return err // TODO
		}
	} else {
		if err = goth_fiber.Logout(c); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error()) // TODO
		}
	}

	return req.RedirectToRoot()
}
