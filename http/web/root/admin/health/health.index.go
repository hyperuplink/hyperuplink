package health

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logichealth "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/health"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type Issue struct {
	LegendKey string
	TextKey   string
	FixClass  string
	FixURL    string
}

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminHealth")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	logicIssues, err := logichealth.Issues(r.Runtime)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var issues []Issue
	for _, issue := range logicIssues {
		issues = append(issues, Issue{
			LegendKey: issue.LegendKey,
			TextKey:   issue.TextKey,
			FixClass:  issue.FixClass,
			FixURL:    route.For(issue.FixRouteID).AsURL(),
		})
	}

	req.SetData("issues", issues)

	return req.Respond()
}
