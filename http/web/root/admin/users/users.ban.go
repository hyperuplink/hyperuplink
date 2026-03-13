package users

import (
	"reflect"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/user"
	"github.com/mrusme/hyperuplink/services/repositories/common"
)

type UserBanForm struct {
	ID string `form:"id" validate:"required,uuid"`
}

func (r *Route) Ban(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminUsers")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(UserBanForm)

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

	usr.SetBannedAt(time.Now())

	err = r.Runtime.Repositories.User.Update(usr)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
