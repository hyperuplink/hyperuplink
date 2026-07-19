package newpost

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicnewpost "xn--gckvb8fzb.com/hyperuplink/logic/root/newpost"
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

	var view *logicnewpost.FormView
	view, err = logicnewpost.View(r.Runtime, req.Perms(), &logicnewpost.FormViewInput{
		ForumSlug: c.Query("forum"),
		TopicSlug: c.Query("topic"),
		ReplyID:   c.Query("reply"),
	})
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.Respond(view)
}
