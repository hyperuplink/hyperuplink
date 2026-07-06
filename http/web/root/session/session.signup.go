package session

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"

	"xn--gckvb8fzb.com/hyperuplink/errs"
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
	Email          string `form:"email" validate:"required,max=254"`
	EmailIsJID     bool   `form:"email_is_jid"`
	Password       string `form:"password" validate:"required,min=8,max=64"`
	PasswordRepeat string `form:"password_repeat" validate:"required,eqcsfield=Password"`
}

func isValidJID(jid string) bool {
	at := strings.IndexByte(jid, '@')
	if at <= 0 || at == len(jid)-1 {
		return false
	}
	if strings.Count(jid, "@") != 1 {
		return false
	}
	if strings.ContainsAny(jid, " \t\r\n") {
		return false
	}
	return true
}

func (r *Route) SignUpShow(c fiber.Ctx) (err error) {
	myRoute := route.For("SessionSignup")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL(),
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(user.GuestRole); ret {
		return rerr
	}

	settingAuth, err := repoSetting.GetByID[setting.Auth](r.Runtime.Repositories.Setting, "auth")
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

	settingAuth, err := repoSetting.GetByID[setting.Auth](r.Runtime.Repositories.Setting, "auth")
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	req.SetData("setting_auth", &settingAuth.JSONValue)

	frm := new(SignUpForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.Respond()
	}

	r.Runtime.Debug("form", frm)

	var isJID bool
	switch settingAuth.JSONValue.AddressType {
	case setting.JID:
		isJID = true
	case setting.EmailAndJID:
		isJID = frm.EmailIsJID
	default:
		isJID = false
	}

	if isJID {
		if !isValidJID(frm.Email) {
			req.Flash.SetErrorsMap(map[string]error{
				"email": fmt.Errorf("%w_email_jid", errs.ErrValidation),
			})
			req.Form.Set(frm)
			return req.Respond()
		}
	} else if verr := validator.New().Var(frm.Email, "email"); verr != nil {
		req.Flash.SetErrorsMap(map[string]error{
			"email": fmt.Errorf("%w_email_email", errs.ErrValidation),
		})
		req.Form.Set(frm)
		return req.Respond()
	}

	// Sign up user
	usr := new(user.User)
	usr.Username = frm.Username
	usr.Role = "user"
	usr.Email = frm.Email
	usr.SetEmailForConfirmation(frm.Email)
	if isJID {
		usr.SetEmailIsJID()
	}
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

	settingCommsEmail, err := repoSetting.GetByID[setting.CommsEmail](r.Runtime.Repositories.Setting, "comms_email")
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	settingCommsXMPP, err := repoSetting.GetByID[setting.CommsXMPP](r.Runtime.Repositories.Setting, "comms_xmpp")
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	// Send confirmation
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

	targetID := settingCommsEmail.JSONValue.TargetID
	if isJID {
		targetID = settingCommsXMPP.JSONValue.TargetID
	}

	err = r.Runtime.Dispatch.SignupConfirmation(
		targetID,
		sc,
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	if usr.Role == user.AdminRole {
		req.Flash.SetInfo(req.In.Ts("signup_success_admin") + " " + usr.EmailConfirmationToken)
	} else if isJID {
		req.Flash.SetInfo(req.In.Ts("signup_success_xmpp"))
	} else {
		req.Flash.SetInfo(req.In.Ts("signup_success"))
	}
	return req.RedirectToRouteID("SessionConfirm")
}
