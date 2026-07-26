package confirm

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type Route struct {
	route.RouteController
}

func New(
	rt *runtime.Runtime,
	router fiber.Router,
) (*Route, error) {
	r := new(Route)

	r.Runtime = rt
	r.Router = router
	r.Path = route.For("SessionConfirm").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("/",
			r.ConfirmShow).Name("confirm.show")
		base.Post("/",
			r.ConfirmCreate).Name("confirm.create")

		base.Get("/"+route.For("SessionConfirmResend").Pathname(),
			r.ConfirmResendShow).Name("confirm.resend.show")
		base.Post("/"+route.For("SessionConfirmResend").Pathname(),
			r.ConfirmResendCreate).Name("confirm.resend.create")
	}, r.Path+".")

	return r, nil
}

func (r *Route) GetRuntime() *runtime.Runtime {
	return r.Runtime
}

func (r *Route) GetPath() string {
	return r.Path
}

func (r *Route) GetEnv() *route.Environment {
	return r.Env
}

type ConfirmForm struct {
	Username string `form:"username" validate:"required,min=2,max=32"`
	Token    string `form:"token" validate:"required,min=8,max=8"`
}

func (r *Route) ConfirmShow(c fiber.Ctx) (err error) {
	myRoute := route.For("SessionConfirm")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL(),
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.GuestRole,
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	return req.Respond()
}

func (r *Route) ConfirmCreate(c fiber.Ctx) (err error) {
	myRoute := route.For("SessionConfirm")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL(),
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.GuestRole,
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(ConfirmForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.Respond()
	}

	r.Runtime.Debug("form", frm)

	var usr *user.User
	usr, err = gh.Repositories(r.Runtime).User.GetByUsername(
		frm.Username,
		common.QueryOptions{
			WithBanned:  false,
			WithSpammed: false,
			WithDeleted: false,
			Limit:       1,
		},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	if usr.EmailConfirmationToken == "" {
		// TODO: Huh?
		if usr.EmailConfirmationSentAt.Valid == true &&
			usr.EmailConfirmationSentAt.Time.IsZero() == false {
			// TODO: Huh??!!?
		}
		// TOOD: Huh...
	}

	if usr.EmailConfirmationToken != frm.Token {
		req.Form.Set(frm)
		return req.RespondError(errs.ErrEmailConfirmationTokenWrong)
	}

	if err := req.Session.Set("local", usr.ID.String()); err != nil {
		req.Session.Reset()
		return req.RespondError(err)
	}

	err = usr.ConfirmEmail(frm.Token)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	err = gh.Repositories(r.Runtime).User.Update(usr)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRouteID("SessionSignin")
}
