package vactivity

import (
	"encoding/json"
	"fmt"
	"strings"

	"xn--gckvb8fzb.com/hyperuplink/models/activity"
)

type VActivity struct {
	activity.Activity

	ActorUsername string `json:"actor_username"`
}

func (m VActivity) Ctx() (ctx activity.Context) {
	if len(m.Context) == 0 {
		return ctx
	}

	if err := json.Unmarshal(m.Context, &ctx); err != nil {
		return activity.Context{}
	}

	return ctx
}

func (m VActivity) ActionKey() string {
	switch m.Kind {
	case activity.AdminVisit:
		return "logs_action_visit"
	case activity.AdminSettingsUpdate:
		return "logs_action_settings_update"
	default:
		return ""
	}
}

func (m VActivity) Details() string {
	ctx := m.Ctx()

	if ctx.Setting == "" {
		return ctx.Path
	}

	if len(ctx.Changed) == 0 {
		return ctx.Setting
	}

	return fmt.Sprintf("%s: %s", ctx.Setting, strings.Join(ctx.Changed, ", "))
}
