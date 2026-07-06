package categories

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/category"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type CategoryMoveForm struct {
	ID string `form:"id" validate:"required,uuid"`
}

func (r *Route) MoveUp(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardCategories")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(CategoryMoveForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", frm)

	cat := new(category.Category)
	cat.ID, err = uuid.Parse(frm.ID)

	err = r.Runtime.Repositories.Category.MoveUp(cat)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}

func (r *Route) MoveDown(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardCategories")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(CategoryMoveForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", frm)

	cat := new(category.Category)
	cat.ID, err = uuid.Parse(frm.ID)

	err = r.Runtime.Repositories.Category.MoveDown(cat)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
