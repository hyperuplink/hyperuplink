package session

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/user"
	goth_fiber "github.com/shareed2k/goth_fiber/v2"
)

func (r *Route) SignOutShow(c fiber.Ctx) (err error) {
	myRoute := route.For("SessionSignout")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL(),
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.UserRole, user.AdminRole,
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
