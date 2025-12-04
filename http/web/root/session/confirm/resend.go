package confirm

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/http/web/request/bcn"
	"github.com/mrusme/hyperuplink/models/user"
)

type ConfirmResendForm struct {
	Email string `form:"email" validate:"required,email"`
}

func (r *Route) ConfirmResendShow(c fiber.Ctx) (err error) {
	req := request.New(r, c,
		[]string{"base"}, route.For("SessionConfirmResend").AsURL(),
		"resend_confirmation")
	req.BCN.Prepend(*bcn.NewBreadcrumb(false, req.In.T("confirm_account"), req.In.T("confirm_account"), req.HrefTo(route.For("SessionConfirm").AsURL())))

	return req.Respond()
}

func (r *Route) ConfirmResendCreate(c fiber.Ctx) (err error) {
	req := request.New(r, c,
		[]string{"base"}, route.For("SessionConfirmResend").AsURL(),
		"resend_confirmation")

	frm := new(ConfirmResendForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.Respond()
	}

	r.Runtime.Debug("form", frm)

	var usr *user.User
	usr, err = r.Runtime.Repositories.User.GetByEmail(frm.Email)
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

	if err := req.Session.Set("local", usr.ID.String()); err != nil {
		req.Session.Reset()
		return req.RespondError(err)
	}

	err = r.Runtime.Repositories.User.Update(usr)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRouteID("SessionSignin")
}
