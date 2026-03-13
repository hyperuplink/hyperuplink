package password

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/errs"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/user"
)

type PasswordUpdateForm struct {
	CurrentPassword   string `form:"current_password" validate:"required,min=8,max=64"`
	NewPassword       string `form:"new_password" validate:"required,min=8,max=64"`
	NewPasswordRepeat string `form:"new_password_repeat" validate:"required,eqcsfield=NewPassword"`
}

func (r *Route) Update(c fiber.Ctx) (err error) {
	myRoute := route.For("AccountPassword")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(PasswordUpdateForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", frm)

	var usr *user.User
	usr, err = req.GetUser()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var match bool
	if match, _, err = usr.CheckPassword(frm.CurrentPassword); !match {
		if err == nil {
			err = errs.ErrPasswordWrong
		}
	}
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	err = usr.SetPassword(frm.NewPassword)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	err = r.Runtime.Repositories.User.Update(usr)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.Flash.SetInfo(req.In.Ts("password_updated"))
	return req.RedirectToRoute(myRoute)
}
