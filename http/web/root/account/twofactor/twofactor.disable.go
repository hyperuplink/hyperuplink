package twofactor

import (
	"errors"
	"reflect"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logictwofactor "xn--gckvb8fzb.com/hyperuplink/logic/root/account/twofactor"
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

	err = logictwofactor.Disable(r.Runtime, usr, frm.CurrentPassword)
	if errors.Is(err, errs.ErrPasswordWrong) {
		req.Flash.SetError(err)
		return req.RedirectToRoute(myRoute)
	}
	if ret, rerr := req.RedirectToRouteOnError(err, myRoute); ret == true {
		return rerr
	}

	req.Session.ClearPendingOTPURL()

	req.Flash.SetInfo(req.In.Ts("twofactor_disabled"))
	return req.RedirectToRoute(myRoute)
}
