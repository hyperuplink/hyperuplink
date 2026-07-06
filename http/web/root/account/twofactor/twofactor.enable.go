package twofactor

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type TwofactorEnableForm struct {
	OTPCode string `form:"otp_code" validate:"required,numeric,len=6"`
}

func (r *Route) Enable(c fiber.Ctx) (err error) {
	myRoute := route.For("AccountTwofactor")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	var usr *user.User
	usr, err = req.GetUser()
	if ret, rerr := req.RedirectToRouteOnError(err, myRoute); ret == true {
		return rerr
	}

	if usr.OTPEnabled {
		return req.RedirectToRoute(myRoute)
	}

	frm := new(TwofactorEnableForm)
	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	pendingURL, ok := req.Session.GetPendingOTPURL()
	if !ok {
		req.Flash.SetError(errs.ErrOTPSetupExpired)
		return req.RedirectToRoute(myRoute)
	}

	var secret string
	secret, err = user.OTPSecretFromURL(pendingURL)
	if ret, rerr := req.RedirectToRouteOnError(err, myRoute); ret == true {
		return rerr
	}

	if !user.ValidateOTP(secret, frm.OTPCode) {
		req.Flash.SetError(errs.ErrOTPCodeWrong)
		return req.RedirectToRoute(myRoute)
	}

	usr.EnableOTP(secret)

	err = r.Runtime.Repositories.User.Update(usr)
	if ret, rerr := req.RedirectToRouteOnError(err, myRoute); ret == true {
		return rerr
	}

	req.Session.ClearPendingOTPURL()

	req.Flash.SetInfo(req.In.Ts("twofactor_enabled"))
	return req.RedirectToRoute(myRoute)
}
