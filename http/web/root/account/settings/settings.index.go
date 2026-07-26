package settings

import (
	"slices"

	"github.com/gofiber/fiber/v3"
	"github.com/zlasd/tzloc"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("AccountSettings")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	var usr *user.User
	usr, err = req.GetUser()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("user", usr)

	timezones := tzloc.GetLocationList()
	slices.Sort(timezones)
	req.SetData("timezones", timezones)

	return req.Respond()
}
