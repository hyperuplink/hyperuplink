package session

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicsession "xn--gckvb8fzb.com/hyperuplink/logic/root/session"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type SignInForm struct {
	// TODO: https://github.com/go-playground/validator/issues/807
	Username string `form:"username" validate:"required,min=2,max=32"`
	Password string `form:"password" validate:"required,min=8,max=64"`
}

func (r *Route) SignInShow(c fiber.Ctx) (err error) {
	myRoute := route.For("SessionSignin")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL(),
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(user.GuestRole); ret {
		return rerr
	}

	return req.Respond()
}

func (r *Route) SignInCreate(c fiber.Ctx) (err error) {
	myRoute := route.For("SessionSignin")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL(),
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.GuestRole,
	); ret {
		return rerr
	}

	frm := new(SignInForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.Respond()
	}

	r.Runtime.Debug("form", frm)

	var usr *user.User
	usr, err = logicsession.CheckPassword(r.Runtime, frm.Username, frm.Password)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	if usr.OTPEnabled {
		req.Session.SetPending2FA("local", usr.ID.String())
		return req.RedirectToRoute(route.For("SessionTwofactor"))
	}

	if err := req.Session.Set("local", usr.ID.String()); err != nil {
		req.Session.Reset()
		return req.RespondError(err)
	}

	return req.RedirectToRoot()
}
