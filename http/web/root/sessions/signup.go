package sessions

import (
	"reflect"

	"github.com/gofiber/fiber/v3"

	"github.com/mrusme/hyperuplink/http/web/modules/errors"
	"github.com/mrusme/hyperuplink/http/web/modules/site"
)

type SignUpForm struct {
	// TODO: https://github.com/go-playground/validator/issues/807
	Username       string `form:"username" validate:"required,min=2,max=32"`
	Email          string `form:"email" validate:"required,email"`
	Password       string `form:"password" validate:"required,min=8,max=64"`
	PasswordRepeat string `form:"password_repeat" validate:"required,eqcsfield=Password"`
}

func (r *Route) SignUpShow(c fiber.Ctx) error {
	sit := site.New(r, c)

	return c.Render("views/session/signup", fiber.Map{
		"Site": sit,
	}, "views/layouts/base")
}

func (r *Route) SignUpCreate(c fiber.Ctx) error {
	sit := site.New(r, c)
	ers := errors.New()
	frm := new(SignUpForm)

	if errmap, ok := sit.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		ers.SetMap(errmap)

		return c.Render("views/session/signup", fiber.Map{
			"Site":   sit,
			"Errors": ers,
		}, "views/layouts/base")
	}

	r.Runtime.Debug("form", frm)

	// TODO: Validate form
	// TODO: Sign up user
	// TODO: Redirect to user profile settings

	return c.Render("views/session/signup", fiber.Map{
		"Site": sit,
	}, "views/layouts/base")
}
