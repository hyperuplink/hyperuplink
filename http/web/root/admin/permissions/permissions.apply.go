package permissions

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/permission"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type PermissionApplyForm struct {
	GroupID    string `form:"group_id" validate:"omitempty,slug,max=32"`
	CategoryID string `form:"category_id" validate:"omitempty,uuid"`
	Level      string `form:"level" validate:"required,oneof=none read read_write read_write_moderate"`
}

func (r *Route) Apply(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminPermissions")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(PermissionApplyForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	isDefault := frm.GroupID == "" && frm.CategoryID == ""
	isGroupMapping := frm.GroupID != "" && frm.CategoryID != ""

	if isDefault == false && isGroupMapping == false {
		req.Flash.SetError(errs.ErrValidation)
		return req.RedirectToRoute(myRoute)
	}

	if isGroupMapping && frm.Level == permission.LevelNone {
		req.Flash.SetError(errs.ErrValidation)
		return req.RedirectToRoute(myRoute)
	}

	level, ok := permission.LevelFromString(frm.Level)
	if ok == false {
		req.Flash.SetError(errs.ErrValidation)
		return req.RedirectToRoute(myRoute)
	}

	var groupID pgtype.Text
	if frm.GroupID != "" {
		groupID = pgtype.Text{String: frm.GroupID, Valid: true}
	}

	var categoryID pgtype.UUID
	if frm.CategoryID != "" {
		uid, perr := uuid.Parse(frm.CategoryID)
		if ret, rerr := req.RespondOnError(perr); ret == true {
			return rerr
		}
		categoryID = pgtype.UUID{Bytes: [16]byte(uid), Valid: true}
	}

	err = r.Runtime.Repositories.Permission.Apply(groupID, categoryID, level)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
