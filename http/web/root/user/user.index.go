package user

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicuser "xn--gckvb8fzb.com/hyperuplink/logic/root/user"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("User")
	req := request.New(r, c, myRoute,
		[]string{"base"}, "user/index",
		"")

	if ret, rerr := req.AccessControl(
		user.GuestRole, // TODO: Maybe remove!
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	var username string = c.Params("user")
	req.UpdateTitle("~" + username)

	var view *logicuser.View
	view, err = logicuser.Show(r.Runtime, username, req.Perms())
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("user", view.User)
	req.SetData("groups", view.Groups)
	req.SetData("topics", &view.Topics)
	req.SetData("replies", &view.Replies)

	return req.Respond()
}
