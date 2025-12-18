package categories

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/category"
	"github.com/mrusme/hyperuplink/models/user"
)

type CategoryDestroyForm struct {
	ID string `form:"id" validate:"required,uuid"`
}

func (r *Route) Destroy(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardCategories")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.GuestRole, // TODO: Remove!
		user.UserRole,  // TODO: Remove!
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(CategoryDestroyForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", frm)

	cat := new(category.Category)
	cat.ID, err = uuid.Parse(frm.ID)

	err = r.Runtime.Repositories.Category.Delete(cat)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
