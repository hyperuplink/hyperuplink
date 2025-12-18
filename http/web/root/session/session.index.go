package session

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/user"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("Session")
	req := request.New(r, c, myRoute,
		[]string{"base"}, "",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.GuestRole,
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	var usr *user.User
	usr, err = req.GetUser()
	if ret, rerr := req.RedirectToRootOnError(err); ret == true {
		return rerr
	}

	req.SetData("user", usr)

	return req.RedirectToRoot()
}
