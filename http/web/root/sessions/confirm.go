package sessions

import (
	"errors"
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/user"
)

type ConfirmForm struct {
	Username string `form:"username" validate:"required,min=2,max=32"`
	Token    string `form:"token" validate:"required,min=8,max=8"`
}

func (r *Route) ConfirmShow(c fiber.Ctx) (err error) {
	req := request.New(r, c, []string{"base"}, "session/confirm",
		"confirm_account")

	return req.Respond()
}

func (r *Route) ConfirmCreate(c fiber.Ctx) (err error) {
	req := request.New(r, c, []string{"base"}, "session/confirm",
		"confirm_account")
	frm := new(ConfirmForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.Respond()
	}

	r.Runtime.Debug("form", frm)

	var usr *user.User
	usr, err = r.Runtime.Repositories.User.GetByUsername(frm.Username)
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
		return req.RespondError(errors.New("email_confirmation_token_wrong"))
	}

	if err := req.Session.Set("local", usr.ID.String()); err != nil {
		req.Session.Reset()
		return req.RespondError(err)
	}

	err = usr.ConfirmEmail(frm.Token)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	err = r.Runtime.Repositories.User.Update(usr)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectTo("sessions/signin")
}
