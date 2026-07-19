package attachments

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicattachments "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/attachments"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Update(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(logicattachments.UpdateInput)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	actorID := uuid.NullUUID{}
	actorID.UUID, actorID.Valid = req.UserUUID()

	err = logicattachments.Update(r.Runtime, actorID, in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
