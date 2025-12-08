package session

import (
	"reflect"

	"github.com/gofiber/fiber/v3"

	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/asyncjob/confirmation/signupconfirmation"
	"github.com/mrusme/hyperuplink/models/setting"
	"github.com/mrusme/hyperuplink/models/user"
	repoSetting "github.com/mrusme/hyperuplink/services/repositories/setting"
)

type SignUpForm struct {
	// TODO: https://github.com/go-playground/validator/issues/807
	Username       string `form:"username" validate:"required,min=2,max=32"`
	Email          string `form:"email" validate:"required,email"`
	Password       string `form:"password" validate:"required,min=8,max=64"`
	PasswordRepeat string `form:"password_repeat" validate:"required,eqcsfield=Password"`
}

func (r *Route) SignUpShow(c fiber.Ctx) (err error) {
	myRoute := route.For("SessionSignup")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL(),
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(user.GuestRole); ret {
		return rerr
	}

	return req.Respond()
}

func (r *Route) SignUpCreate(c fiber.Ctx) (err error) {
	myRoute := route.For("SessionSignup")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL(),
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(user.GuestRole); ret {
		return rerr
	}

	frm := new(SignUpForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.Respond()
	}

	r.Runtime.Debug("form", frm)

	// Sign up user
	usr := new(user.User)
	usr.Username = frm.Username
	usr.Role = "user"
	usr.Email = frm.Email
	usr.SetEmailForConfirmation(frm.Email)
	usr.Language = req.In.Lang()

	err = usr.SetPassword(frm.Password)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	usr.ID, err = r.Runtime.Repositories.User.Create(usr)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	settingSystem, err := repoSetting.GetByID[setting.System](r.Runtime.Repositories.Setting, "system")
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	// Send confirmation email
	sc, err := signupconfirmation.New(
		&settingSystem.JSONValue,
		usr,
		req.In.T("signup_confirmation_subject"),
		usr.EmailConfirmationToken,
		myRoute.AsURL(),
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	err = r.Runtime.Dispatch.SignupConfirmation(
		"notifications", // TODO: Replace with System.EmailTarget
		sc,
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.Flash.SetInfo(req.In.T("signup_success"))
	return req.RedirectToRouteID("SessionConfirm")
}
