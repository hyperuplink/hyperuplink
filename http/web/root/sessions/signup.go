package sessions

import (
	"reflect"

	"github.com/gofiber/fiber/v3"

	"github.com/mrusme/hyperuplink/http/web/request"
)

type SignUpForm struct {
	// TODO: https://github.com/go-playground/validator/issues/807
	Username       string `form:"username" validate:"required,min=2,max=32"`
	Email          string `form:"email" validate:"required,email"`
	Password       string `form:"password" validate:"required,min=8,max=64"`
	PasswordRepeat string `form:"password_repeat" validate:"required,eqcsfield=Password"`
}

func (r *Route) SignUpShow(c fiber.Ctx) error {
	req := request.New(r, c)

	return req.Respond("base", "session/signup")
}

func (r *Route) SignUpCreate(c fiber.Ctx) error {
	req := request.New(r, c)
	frm := new(SignUpForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.Respond("base", "session/signup")
	}

	r.Runtime.Debug("form", frm)

	// TODO: Validate form
	// TODO: Sign up user
	// TODO: Redirect to user profile settings

	return req.Respond("base", "session/signup")
}
