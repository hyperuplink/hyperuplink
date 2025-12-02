package sessions

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/web/request"
	goth_fiber "github.com/shareed2k/goth_fiber/v2"
)

func (r *Route) SignOutShow(c fiber.Ctx) error {
	req := request.New(r, c)

	if req.Session.GetProvider() == "local" {
		if err := req.Session.Reset(); err != nil {
			return err // TODO
		}
	} else {
		if err := goth_fiber.Logout(c); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error()) // TODO
		}
	}

	return req.RedirectToRoot()
}
