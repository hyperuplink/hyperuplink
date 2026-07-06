package users

import (
	"reflect"
	"time"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

type UserUnbanForm struct {
	ID string `form:"id" validate:"required,uuid"`
}

func (r *Route) Unban(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminUsers")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(UserUnbanForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", frm)

	usr, err := r.Runtime.Repositories.User.GetByID(
		frm.ID,
		common.QueryOptions{
			WithBanned:  true,
			WithDeleted: true,
			Limit:       1,
		},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	usr.SetBannedAt(time.Time{})

	err = r.Runtime.Repositories.User.Update(usr)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
