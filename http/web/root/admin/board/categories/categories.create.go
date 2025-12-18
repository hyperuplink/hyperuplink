package categories

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/category"
	"github.com/mrusme/hyperuplink/models/user"
)

type CategoryCreateForm struct {
	Name string `form:"name" validate:"required,min=1,max=32"`
	Slug string `form:"slug" validate:"required,slug,min=1,max=32"`
}

func (r *Route) Create(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardCategories")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(CategoryCreateForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", frm)

	cat := new(category.Category)
	cat.Name = frm.Name
	cat.Slug = frm.Slug

	_, err = r.Runtime.Repositories.Category.Create(cat)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
