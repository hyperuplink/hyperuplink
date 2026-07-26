package password

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicpassword "xn--gckvb8fzb.com/hyperuplink/logic/root/account/password"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

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

	in := new(logicpassword.UpdateInput)

	if ok := req.ValidateForm(in, reflect.TypeOf(*in)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	var usr *user.User
	usr, err = req.GetUser()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("user", usr)

	err = logicpassword.Update(r.Runtime, usr, in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.Flash.SetInfo(req.In.Ts("password_updated"))
	return req.RedirectToRoute(myRoute)
}
