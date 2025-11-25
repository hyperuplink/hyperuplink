package sessions

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/mrusme/hyperuplink/http/web/site"
)

type SignUpForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func (r *Route) SignUpShow(c fiber.Ctx) error {
	return c.Render("views/session/signin", fiber.Map{
		"Site": site.New(r, c),
	}, "views/layouts/base")
}

func (r *Route) SignUpCreate(c fiber.Ctx) error {
	s := site.New(r, c)
	sess := session.FromContext(c)
	f := new(SignUpForm)

	if err := c.Bind().Form(f); err != nil {
		return err // TODO
	}

	r.Runtime.Debug("form", f)

	// TODO: Validate form
	// TODO: Sign up user
	// TODO: Redirect to user profile settings

	return c.Render("views/session/signin", fiber.Map{
		"Site": s,
	}, "views/layouts/base")
}
