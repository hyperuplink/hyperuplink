package attachments

import (
	"errors"
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicattachments "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/attachments"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Update(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardAttachments")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(logicattachments.UpdateInput)

	if ok := req.ValidateForm(in, reflect.TypeOf(*in)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", in)

	actorID := uuid.NullUUID{}
	actorID.UUID, actorID.Valid = req.Session.GetUserUUID()

	err = logicattachments.Update(r.Runtime, actorID, in)
	if errors.Is(err, errs.ErrInvalidUploadFormat) ||
		errors.Is(err, errs.ErrInvalidStorageProvider) {
		req.Flash.SetError(err)
		return req.RedirectToRoute(myRoute)
	}
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
