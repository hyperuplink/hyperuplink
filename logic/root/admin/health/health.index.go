package health

import (
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/config"
)

type Issue struct {
	LegendKey  string `json:"legend_key"`
	TextKey    string `json:"text_key"`
	FixClass   string `json:"fix_class"`
	FixRouteID string `json:"fix_route_id"`
}

func Issues(rt *runtime.Runtime) (issues []Issue, err error) {
	targets, err := rt.Config().Targets()
	if err != nil {
		return nil, err
	}

	if issue, found := checkXMPPInsecureSkipVerify(rt, targets); found {
		issues = append(issues, issue)
	}

	return issues, nil
}

func checkXMPPInsecureSkipVerify(
	rt *runtime.Runtime,
	targets config.Targets,
) (issue Issue, found bool) {
	if rt.IsDevelopmentMode() {
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
			LegendKey:  "health_xmpp_insecure_skip_verify",
			TextKey:    "health_xmpp_insecure_skip_verify_text",
			FixClass:   "warn",
			FixRouteID: "AdminCommsXmpp",
		}, true
	}

	return issue, false
}
