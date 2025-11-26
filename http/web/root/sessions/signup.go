package sessions

import (
	"reflect"

	"github.com/gofiber/fiber/v3"

	// "github.com/gofiber/fiber/v3/middleware/session"
	"github.com/mrusme/hyperuplink/http/web/site"
)

type SignUpForm struct {
	// TODO: https://github.com/go-playground/validator/issues/807
	Username       string `form:"username" validate:"required,min=2,max=32"`
	Email          string `form:"email" validate:"required,email"`
	Password       string `form:"password" validate:"required,min=8,max=64"`
	PasswordRepeat string `form:"password_repeat" validate:"required,eqcsfield=Password"`
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

	if errmap, ok := s.ValidateForm(f, reflect.TypeOf(*f)); !ok {
		return c.Render("views/session/signup", fiber.Map{
			"Site":   s,
			"Errors": errmap,
		}, "views/layouts/base")
	}

	r.Runtime.Debug("form", f)

	// TODO: Validate form
	// TODO: Sign up user
	// TODO: Redirect to user profile settings

	return c.Render("views/session/signup", fiber.Map{
		"Site": s,
	}, "views/layouts/base")
}
