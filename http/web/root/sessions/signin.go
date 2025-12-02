package sessions

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/web/request"
)

type SignInForm struct {
	// TODO: https://github.com/go-playground/validator/issues/807
	Username string `form:"username" validate:"required,min=2,max=32"`
	Password string `form:"password" validate:"required,min=8,max=64"`
}

func (r *Route) SignInShow(c fiber.Ctx) error {
	req := request.New(r, c)

	return req.Respond("base", "session/signin")
}

func (r *Route) SignInCreate(c fiber.Ctx) error {
	req := request.New(r, c)
	frm := new(SignInForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.Respond("base", "session/signin")
	}

	r.Runtime.Debug("form", frm)

	if frm.Username == "user" && frm.Password == "pass" {
		if err := req.Session.Set("local", "2941476f-2ae0-4c3e-a459-1ef5d8dd6ca9"); err != nil {
			req.Session.Reset()
		}

		return req.RedirectToRoot()
	}

	return req.Respond("base", "session/signin")
}
