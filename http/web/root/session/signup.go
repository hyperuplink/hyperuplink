package session

import (
	"reflect"

	"github.com/gofiber/fiber/v3"

	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/user"
)

type SignUpForm struct {
	// TODO: https://github.com/go-playground/validator/issues/807
	Username       string `form:"username" validate:"required,min=2,max=32"`
	Email          string `form:"email" validate:"required,email"`
	Password       string `form:"password" validate:"required,min=8,max=64"`
	PasswordRepeat string `form:"password_repeat" validate:"required,eqcsfield=Password"`
}

func (r *Route) SignUpShow(c fiber.Ctx) (err error) {
	req := request.New(r, c,
		[]string{"base"}, route.For("SessionSignup").AsURL(),
		"sign_up_noun")
	if ret, rerr := req.AccessControl(user.GuestRole); ret {
		return rerr
	}

	return req.Respond()
}

func (r *Route) SignUpCreate(c fiber.Ctx) (err error) {
	req := request.New(r, c,
		[]string{"base"}, route.For("SessionSignup").AsURL(),
		"sign_up_noun")
	if ret, rerr := req.AccessControl(user.GuestRole); ret {
		return rerr
	}

	frm := new(SignUpForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.Respond()
	}

	r.Runtime.Debug("form", frm)

	// TODO: Sign up user
	usr := new(user.User)
	usr.Username = frm.Username
	usr.Role = "user"
	usr.Email = frm.Email
	usr.SetEmailForConfirmation(frm.Email)

	err = usr.SetPassword(frm.Password)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	usr.ID, err = r.Runtime.Repositories.User.Create(usr)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	// TODO: Trigger confirmation mail

	return req.RedirectToRouteID("SessionConfirm")
}
