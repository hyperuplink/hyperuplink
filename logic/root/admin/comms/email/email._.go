package email

import (
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/services/config"
)

type UpdateInput struct {
	TargetID string `json:"target_id" form:"target_id" validate:"omitempty,max=64"`
}

type View struct {
	CommsEmail     *setting.CommsEmail `json:"comms_email"`
	EmailTargets   config.Targets      `json:"email_targets"`
	SelectedTarget *config.Target      `json:"selected_target,omitempty"`
}
