package reports

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicreport "xn--gckvb8fzb.com/hyperuplink/logic/root/report"
	"xn--gckvb8fzb.com/hyperuplink/models/postevent"
)

var reasons = map[string]postevent.ReportType{
	"spam":       postevent.Spam,
	"misconduct": postevent.Misconduct,
	"illegal":    postevent.Illegal,
}

// @Summary		Report a post
// @Description	Files a report about a topic or a reply, addressed by UUID or short
// @Description	id, for the administrators to read. One report per viewer and post;
// @Description	a second one is a conflict.
// @Tags			board
// @Accept			json
// @Produce		json
// @Param			request	body		reports.ReportBody	true	"The post and the reason"
// @Success		201		{object}	object{target=string,id=string,reason=string}
// @Failure		401		{object}	problem.Problem
// @Failure		404		{object}	problem.Problem
// @Failure		409		{object}	problem.Problem
// @Failure		422		{object}	problem.Problem
// @Security		BearerAuth
// @Router			/reports [post]
func (r *Route) Create(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireUser(); ret {
		return rerr
	}

	in := new(ReportBody)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logicreport.Create(req.Context(), r.Runtime, req.User().ID, &logicreport.CreateInput{
		Target:     in.Target,
		ID:         in.ID,
		ReportType: int(reasons[in.Reason]),
	})
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondCreated(fiber.Map{
		"target": in.Target,
		"id":     in.ID,
		"reason": in.Reason,
	})
}
