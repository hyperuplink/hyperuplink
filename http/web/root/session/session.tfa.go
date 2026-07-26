package session

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

const maxTfaAttempts = 5

type TfaForm struct {
	OTPCode string `form:"otp_code" validate:"required,numeric,len=6"`
}

func (r *Route) TfaShow(c fiber.Ctx) (err error) {
	myRoute := route.For("SessionTwofactor")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL(),
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(user.GuestRole); ret {
		return rerr
	}

	if _, _, ok := req.Session.GetPending2FA(); !ok {
		return req.RedirectToRoute(route.For("SessionSignin"))
	}

	return req.Respond()
}

func (r *Route) TfaCreate(c fiber.Ctx) (err error) {
	myRoute := route.For("SessionTwofactor")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL(),
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(user.GuestRole); ret {
		return rerr
	}

	userID, provider, ok := req.Session.GetPending2FA()
	if !ok {
		return req.RedirectToRoute(route.For("SessionSignin"))
	}

	frm := new(TfaForm)
	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	var usr *user.User
	usr, err = gh.Repositories(r.Runtime).User.GetByID(
		userID,
		common.QueryOptions{
			WithBanned:  false,
			WithSpammed: false,
			WithDeleted: false,
			Limit:       1,
		},
	)
	if err != nil {
		req.Session.ClearPending2FA()
		req.Flash.SetError(errs.ErrUsernamePasswordWrong)
		return req.RedirectToRoute(route.For("SessionSignin"))
	}

	if usr.OTPEnabled && !user.ValidateOTP(usr.OTPSecret, frm.OTPCode) {
		if req.Session.IncrementPending2FAAttempts() >= maxTfaAttempts {
			req.Session.ClearPending2FA()
			req.Flash.SetError(errs.ErrOTPTooManyAttempts)
			return req.RedirectToRoute(route.For("SessionSignin"))
		}
		req.Flash.SetError(errs.ErrOTPCodeWrong)
		return req.RedirectToRoute(myRoute)
	}

	if err := req.Session.Set(provider, usr.ID.String()); err != nil {
		req.Session.Reset()
		return req.RespondError(err)
	}

	req.Session.ClearPending2FA()

	return req.RedirectToRoot()
}
