package health

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/services/config"
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

	var targets config.Targets
	targets, err = r.Runtime.Config.Targets()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var issues []Issue

	if issue, found := r.checkXMPPInsecureSkipVerify(targets); found {
		issues = append(issues, issue)
	}

	req.SetData("issues", issues)

	return req.Respond()
}

func (r *Route) checkXMPPInsecureSkipVerify(
	targets config.Targets,
) (issue Issue, found bool) {
	if r.Runtime.IsDevelopmentMode() {
		return issue, false
	}

	for _, target := range targets {
		if target.Type != config.TargetTypeXMPP {
			continue
		}
		if !target.XMPP.InsecureSkipVerify {
			continue
		}

		return Issue{
			LegendKey: "health_xmpp_insecure_skip_verify",
			TextKey:   "health_xmpp_insecure_skip_verify_text",
			FixClass:  "warn",
			FixURL:    route.For("AdminCommsXmpp").AsURL(),
		}, true
	}

	return issue, false
}
