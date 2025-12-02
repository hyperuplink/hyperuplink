package sessions

import (
	"reflect"
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/user"
)

type SignUpForm struct {
	// TODO: https://github.com/go-playground/validator/issues/807
	Username       string `form:"username" validate:"required,min=2,max=32"`
	Email          string `form:"email" validate:"required,email"`
	Password       string `form:"password" validate:"required,min=8,max=64"`
	PasswordRepeat string `form:"password_repeat" validate:"required,eqcsfield=Password"`
}

func (r *Route) SignUpShow(c fiber.Ctx) error {
	req := request.New(r, c)

	return req.Respond("base", "session/signup")
}

func (r *Route) SignUpCreate(c fiber.Ctx) error {
	var err error

	req := request.New(r, c)
	frm := new(SignUpForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.Respond("base", "session/signup")
	}

	r.Runtime.Debug("form", frm)

	// TODO: Sign up user
	usr := new(user.User)
	usr.Username = frm.Username
	usr.Role = "user"
	if ret, rerr := req.RespondOnError("base", "session/signup",
		usr.SetPassword(frm.Password)); ret == true {
		return rerr
	}
	usr.Email = frm.Email
	usr.EmailUnconfirmed = frm.Email
	usr.ResetEmailConfirmationToken()
	usr.EmailConfirmationSentAt = time.Now()

	usr.ID, err = r.Runtime.Repositories.User.Create(usr)
	if ret, rerr := req.RespondOnError("base", "session/signup",
		r.Runtime.Database.ConvertError(err)); ret == true {
		return rerr
	}

	// TODO: Trigger confirmation mail

	// TODO: Redirect to user profile settings
	return req.Respond("base", "session/signup")
}
