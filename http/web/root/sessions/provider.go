package sessions

import (
	"github.com/gofiber/fiber/v3"
	goth_fiber "github.com/shareed2k/goth_fiber/v2"
)

func (r *Route) ProviderCallbackShow(c fiber.Ctx) error {
	user, err := goth_fiber.CompleteUserAuth(c, goth_fiber.CompleteUserAuthOptions{
		ShouldLogout: false,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString(user.Email)
}
