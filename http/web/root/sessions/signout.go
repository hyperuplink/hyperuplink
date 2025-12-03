package sessions

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/web/request"
	goth_fiber "github.com/shareed2k/goth_fiber/v2"
)

func (r *Route) SignOutShow(c fiber.Ctx) (err error) {
	req := request.New(r, c, []string{"base"}, "session/signout", "sign_out_noun")
	if ret, rerr := req.AccessControl(
		request.UserRole, request.ModRole, request.AdminRole,
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
