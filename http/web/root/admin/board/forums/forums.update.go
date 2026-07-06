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

type ForumUpdateForm struct {
	ID          string `form:"id" validate:"required,uuid"`
	Name        string `form:"name" validate:"required,min=1,max=32"`
	Slug        string `form:"slug" validate:"required,slug,min=1,max=32"`
	CategoryID  string `form:"category_id" validate:"required,uuid"`
	Description string `form:"description" validate:"min=0,max=128"`
}

func (r *Route) Update(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardForums")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(ForumUpdateForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", frm)

	fum := new(forum.Forum)
	fum.ID, err = uuid.Parse(frm.ID)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	fum.Name = frm.Name
	fum.Slug = frm.Slug
	fum.CategoryID, err = uuid.Parse(frm.CategoryID)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	fum.Description = frm.Description

	err = r.Runtime.Repositories.Forum.Update(fum)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
