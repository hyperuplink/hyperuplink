package permissions

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type PermissionRemoveForm struct {
	GroupID    string `form:"group_id" validate:"required,slug,max=32"`
	CategoryID string `form:"category_id" validate:"required,uuid"`
}

func (r *Route) PermissionRemove(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminPermissions")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(PermissionRemoveForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	uid, perr := uuid.Parse(frm.CategoryID)
	if ret, rerr := req.RespondOnError(perr); ret == true {
		return rerr
	}

	groupID := pgtype.Text{String: frm.GroupID, Valid: true}
	categoryID := pgtype.UUID{Bytes: [16]byte(uid), Valid: true}

	err = r.Runtime.Repositories.Permission.Remove(groupID, categoryID)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
