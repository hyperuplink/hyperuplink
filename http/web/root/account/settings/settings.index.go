package settings

import (
	"slices"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/user"
	"github.com/zlasd/tzloc"
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
