package sessions

import (
	"errors"
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/user"
)

type SignInForm struct {
	// TODO: https://github.com/go-playground/validator/issues/807
	Username string `form:"username" validate:"required,min=2,max=32"`
	Password string `form:"password" validate:"required,min=8,max=64"`
}

func (r *Route) SignInShow(c fiber.Ctx) (err error) {
	req := request.New(r, c, []string{"base"}, "session/signin", "sign_in_noun")

	return req.Respond()
}

func (r *Route) SignInCreate(c fiber.Ctx) (err error) {
	req := request.New(r, c, []string{"base"}, "session/signin", "sign_in_noun")
	frm := new(SignInForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.Respond()
	}

	r.Runtime.Debug("form", frm)

	var usr *user.User
	usr, err = r.Runtime.Repositories.User.GetByUsername(frm.Username)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var match bool = false
	if match, _, err = usr.CheckPassword(frm.Password); !match {
		if err == nil {
			err = errors.New("username_password_wrong")
		}
	}
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	if err := req.Session.Set("local", usr.ID.String()); err != nil {
		req.Session.Reset()
		return req.RespondError(err)
	}

	return req.RedirectToRoot()
}
