package report

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicreport "xn--gckvb8fzb.com/hyperuplink/logic/root/report"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary	Show the post a report would be about
// @Tags		board
// @Produce	json
// @Param		target	query		string	true	"The kind of post, either topic or reply"
// @Param		id		query		string	true	"The identifier of the post"
// @Success	200		{object}	object{target=string,id=string,post_text=string}
// @Failure	401		{object}	request.ErrorResponse
// @Failure	403		{object}	request.ErrorResponse
// @Failure	404		{object}	request.ErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/report [get]
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
