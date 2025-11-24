package sessions

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/web/site"
)

type SignInForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func (r *Route) SignInShow(c fiber.Ctx) error {
	return c.Render("views/session/signin", fiber.Map{
		"Site": site.New(r, c),
	}, "views/layouts/base")
}

func (r *Route) SignInCreate(c fiber.Ctx) error {
	f := new(SignInForm)

	if err := c.Bind().Form(f); err != nil {
		return err
	}

	r.Runtime.Debug("form", f)

	return c.Render("views/session/signin", fiber.Map{
		"Site": site.New(r, c),
	}, "views/layouts/base")
}
