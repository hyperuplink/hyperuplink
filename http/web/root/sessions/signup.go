package sessions

import (
	"github.com/gofiber/fiber/v3"
	// "github.com/gofiber/fiber/v3/middleware/session"
	"github.com/mrusme/hyperuplink/http/web/site"
)

type SignUpForm struct {
	// TODO: https://github.com/go-playground/validator/issues/807
	Username string `form:"username" validate:"required,min=2,max=32"`
	Password string `form:"password" validate:"required,min=8,max=64"`
}

func (r *Route) SignUpShow(c fiber.Ctx) error {
	s := site.New(r, c)

	return c.Render("views/session/signup", fiber.Map{
		"Site": s,
	}, "views/layouts/base")
}

func (r *Route) SignUpCreate(c fiber.Ctx) error {
	s := site.New(r, c)
	// sess := session.FromContext(c)
	f := new(SignUpForm)

	if err := c.Bind().Form(f); err != nil {
		return err // TODO
	}

	r.Runtime.Debug("form", f)

	// TODO: Validate form
	// TODO: Sign up user
	// TODO: Redirect to user profile settings

	return c.Render("views/session/signup", fiber.Map{
		"Site": s,
	}, "views/layouts/base")
}
