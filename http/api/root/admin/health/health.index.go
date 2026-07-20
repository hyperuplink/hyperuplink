package health

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logichealth "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/health"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		List the health issues
// @Description	Each issue names something the board has found wrong with its own
// @Description	configuration, so an empty list is the healthy answer.
// @Tags			admin
// @Produce		json
// @Success		200	{object}	object{issues=[]logichealth.Issue}
// @Failure		401	{object}	request.ErrorResponse
// @Failure		403	{object}	request.ErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/admin/health [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	issues, err := logichealth.Issues(r.Runtime)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.Respond(fiber.Map{
		"issues": issues,
	})
}
