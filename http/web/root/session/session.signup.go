package session

import (
	"reflect"
	"slices"
	"strings"

	"github.com/gofiber/fiber/v3"

	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/confirmation/signupconfirmation"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	repoSetting "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
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

	if ret, rerr := req.AccessControl(
		user.GuestRole,
	); ret {
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

	promoteAdmin := r.Runtime.Config.UsersPromoteAdmin()
	if slices.Index(promoteAdmin, strings.ToLower(usr.Email)) > -1 {
		usr.Role = user.AdminRole
	}

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
		req.In.Ts("signup_confirmation_subject"),
		usr.EmailConfirmationToken,
		myRoute.AsURL(),
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	r.Runtime.Error("sc", sc)
	r.Runtime.Error("sc", sc)
	r.Runtime.Error("sc", sc)
	r.Runtime.Error("sc", sc)
	r.Runtime.Error("sc", sc)

	err = r.Runtime.Dispatch.SignupConfirmation(
		"notifications", // TODO: Replace with setting.CommsEmail
		sc,
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	if usr.Role == user.AdminRole {
		req.Flash.SetInfo(req.In.Ts("signup_success_admin") + " " + usr.EmailConfirmationToken)
	} else {
		req.Flash.SetInfo(req.In.Ts("signup_success"))
	}
	return req.RedirectToRouteID("SessionConfirm")
}
