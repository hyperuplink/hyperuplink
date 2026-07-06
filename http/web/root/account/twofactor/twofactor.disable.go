package twofactor

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type TwofactorDisableForm struct {
	CurrentPassword string `form:"current_password" validate:"required,min=8,max=64"`
}

func (r *Route) Disable(c fiber.Ctx) (err error) {
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

	if !usr.OTPEnabled {
		return req.RedirectToRoute(myRoute)
	}

	frm := new(TwofactorDisableForm)
	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	var match bool
	if match, _, err = usr.CheckPassword(frm.CurrentPassword); !match {
		if err == nil {
			err = errs.ErrPasswordWrong
		}
	}
	if ret, rerr := req.RedirectToRouteOnError(err, myRoute); ret == true {
		return rerr
	}

	usr.DisableOTP()

	err = r.Runtime.Repositories.User.Update(usr)
	if ret, rerr := req.RedirectToRouteOnError(err, myRoute); ret == true {
		return rerr
	}

	req.Session.ClearPendingOTPURL()

	req.Flash.SetInfo(req.In.Ts("twofactor_disabled"))
	return req.RedirectToRoute(myRoute)
}
