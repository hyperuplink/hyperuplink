package report

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicreport "xn--gckvb8fzb.com/hyperuplink/logic/root/report"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	target := c.Query("target")
	id := c.Query("id")

	post, err := logicreport.ResolvePost(r.Runtime, target, id)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.Respond(fiber.Map{
		"target":    target,
		"id":        id,
		"post_text": post.Text,
	})
}
