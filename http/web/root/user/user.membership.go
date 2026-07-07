package user

import (
	"reflect"
	"slices"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

type MembershipForm struct {
	MemberOf []string `form:"member_of"`
}

func (r *Route) Membership(c fiber.Ctx) (err error) {
	myRoute := route.For("User")
	req := request.New(r, c, myRoute,
		[]string{"base"}, "user/index",
		"")

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	username := c.Params("user")

	frm := new(MembershipForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute.Fill(map[string]string{"user": username}))
	}

	usr, err := r.Runtime.Repositories.User.GetByUsername(
		username,
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

	groups, err := r.Runtime.Repositories.Group.All(common.QueryOptions{})
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	memberOf := []string{}
	for _, grp := range *groups {
		if slices.Contains(frm.MemberOf, grp.ID) {
			memberOf = append(memberOf, grp.ID)
		}
	}
	usr.MemberOf = memberOf

	err = r.Runtime.Repositories.User.Update(usr)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute.Fill(map[string]string{"user": username}))
}
