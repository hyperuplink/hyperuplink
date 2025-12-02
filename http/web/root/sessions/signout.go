package sessions

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/web/modules/session"
	"github.com/mrusme/hyperuplink/http/web/modules/site"
	goth_fiber "github.com/shareed2k/goth_fiber/v2"
)

func (r *Route) SignOutShow(c fiber.Ctx) error {
	sit := site.New(r, c)
	ses := session.New(c)

	if ses.GetProvider() == "local" {
		if err := ses.Reset(); err != nil {
			return err // TODO
		}
	} else {
		if err := goth_fiber.Logout(c); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error()) // TODO
		}
	}

	return c.Redirect().To(sit.GetRelRoot())
}
