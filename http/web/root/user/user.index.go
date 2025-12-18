package user

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/user"
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

	// var fum *vforum.VForum
	// fum, err = r.Runtime.Repositories.Forum.VGetBySlug(
	// 	forum_slug,
	// 	common.QueryOptions{
	// 		Limit: 1,
	// 	},
	// )
	// if ret, rerr := req.RespondOnError(err); ret == true {
	// 	return rerr
	// }
	//
	// req.SetData("forum", fum)

	// req.UpdateTitle(fum.Name)

	return req.Respond()
}
