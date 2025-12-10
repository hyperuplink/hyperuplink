package categories

import (
	"reflect"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/category"
	"github.com/mrusme/hyperuplink/models/user"
)

type CategoryUpdateForm struct {
	ID   string `form:"id" validate:"required,uuid"`
	Name string `form:"name" validate:"required,min=1,max=32"`
	Slug string `form:"slug" validate:"required,slug,min=1,max=32"`
}

func (r *Route) Update(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardCategories")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(user.GuestRole); ret {
		return rerr
	}

	frm := new(CategoryUpdateForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", frm)

	cat := new(category.Category)
	cat.ID, err = uuid.Parse(frm.ID)
	cat.Name = frm.Name
	cat.Slug = frm.Slug
	cat.SetUpdatedAt(time.Now())

	err = r.Runtime.Repositories.Category.Update(cat)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
