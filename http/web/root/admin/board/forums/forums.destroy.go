package forums

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/forum"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type ForumDestroyForm struct {
	ID string `form:"id" validate:"required,uuid"`
}

func (r *Route) Destroy(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardForums")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(ForumDestroyForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", frm)

	fum := new(forum.Forum)
	fum.ID, err = uuid.Parse(frm.ID)

	err = r.Runtime.Repositories.Forum.Delete(fum)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
