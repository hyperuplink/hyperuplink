package forums

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/forum"
	"github.com/mrusme/hyperuplink/models/user"
)

type ForumCreateForm struct {
	Name       string `form:"name" validate:"required,min=1,max=32"`
	Slug       string `form:"slug" validate:"required,slug,min=1,max=32"`
	CategoryID string `form:"category_id" validate:"required,uuid"`
}

func (r *Route) Create(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardForums")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(ForumCreateForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", frm)

	fum := new(forum.Forum)
	fum.Name = frm.Name
	fum.Slug = frm.Slug
	fum.CategoryID, err = uuid.Parse(frm.CategoryID)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	_, err = r.Runtime.Repositories.Forum.Create(fum)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
