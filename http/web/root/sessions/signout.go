package sessions

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/mrusme/hyperuplink/http/web/site"
	goth_fiber "github.com/shareed2k/goth_fiber/v2"
)

func (r *Route) SignOutShow(c fiber.Ctx) error {
	s := site.New(r, c)
	sess := session.FromContext(c)

	if sess.Get("auth") == "local" {
		if err := sess.Reset(); err != nil {
			return err // TODO
		}
	} else {
		if err := goth_fiber.Logout(c); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error()) // TODO
		}
	}

	return c.Redirect().To(s.GetRelRoot())
}
