package session

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v3"

	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicsession "xn--gckvb8fzb.com/hyperuplink/logic/root/session"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	repoSetting "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func (r *Route) SignUpShow(c fiber.Ctx) (err error) {
	myRoute := route.For("SessionSignup")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL(),
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(user.GuestRole); ret {
		return rerr
	}

	settingAuth, err := repoSetting.GetByID[setting.Auth](gh.Repositories(r.Runtime).Setting, "auth")
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	req.SetData("setting_auth", &settingAuth.JSONValue)

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

	settingAuth, err := repoSetting.GetByID[setting.Auth](gh.Repositories(r.Runtime).Setting, "auth")
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	req.SetData("setting_auth", &settingAuth.JSONValue)

	in := new(logicsession.SignUpInput)

	if ok := req.ValidateForm(in, reflect.TypeOf(*in)); !ok {
		return req.Respond()
	}

	if in.PasswordRepeat == "" {
		req.Flash.SetErrorsMap(map[string]error{
			"password_repeat": fmt.Errorf("%w_passwordrepeat_required", errs.ErrValidation),
		})
		req.Form.Set(in)
		return req.Respond()
	}

	in.Language = req.In.Lang()

	r.Runtime.Debug("form", in)

	usr, err := logicsession.SignUp(r.Runtime, in, logicsession.SignUpOptions{Activate: false})
	if err != nil {
		var fe *logicsession.FieldError
		if errors.As(err, &fe) {
			req.Flash.SetErrorsMap(map[string]error{fe.Field: fe.Err})
			req.Form.Set(in)
			return req.Respond()
		}
		if ret, rerr := req.RespondOnError(err); ret == true {
			return rerr
		}
	}

	err = logicsession.SendSignupConfirmation(
		r.Runtime,
		usr,
		req.In.Ts("signup_confirmation_subject"),
		myRoute.AsURL(),
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	if usr.Role == user.AdminRole {
		req.Flash.SetInfo(req.In.Ts("signup_success_admin") + " " + usr.EmailConfirmationToken)
	} else if usr.EmailIsJID {
		req.Flash.SetInfo(req.In.Ts("signup_success_xmpp"))
	} else {
		req.Flash.SetInfo(req.In.Ts("signup_success"))
	}
	return req.RedirectToRouteID("SessionConfirm")
}
