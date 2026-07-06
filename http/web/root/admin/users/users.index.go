package users

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminUsers")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	var users *[]user.User
	users, err = r.Runtime.Repositories.User.All(common.QueryOptions{
		WithBanned:  true,
		WithDeleted: true,
		OrderBy:     "created_at",
		Order:       common.Descending,
	})
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("users", users)

	return req.Respond()
}
