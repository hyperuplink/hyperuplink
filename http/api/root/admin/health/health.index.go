package health

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logichealth "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/health"
)

// @Summary		List the health issues
// @Description	Each issue names something the board has found wrong with its own
// @Description	configuration; an empty list is the healthy answer.
// @Tags			admin
// @Produce		json
// @Success		200	{object}	object{issues=[]logichealth.Issue}
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/health [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	issues, err := logichealth.Issues(req.Context(), r.Runtime)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}
	if issues == nil {
		issues = []logichealth.Issue{}
	}

	return req.Respond(fiber.Map{"issues": issues})
}
