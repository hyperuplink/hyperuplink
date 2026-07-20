package board

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary	List the board administration sections
// @Tags		admin
// @Produce	json
// @Success	200	{object}	object{sections=[]string}
// @Failure	401	{object}	request.ErrorResponse
// @Failure	403	{object}	request.ErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/admin/board [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	return req.Respond(fiber.Map{
		"sections": []string{
			"attachments",
			"categories",
			"forums",
			"profiles",
			"themes",
			"topics",
		},
	})
}
