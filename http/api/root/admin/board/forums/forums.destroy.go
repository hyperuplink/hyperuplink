package forums

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicforums "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/forums"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Destroy(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	err = logicforums.Destroy(r.Runtime, &logicforums.DestroyInput{
		ID: c.Params("id"),
	})
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
