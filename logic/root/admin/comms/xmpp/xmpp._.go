package xmpp

import (
	"xn--gckvb8fzb.com/glides/services/config"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
)

type UpdateInput struct {
	TargetID string `json:"target_id" form:"target_id" validate:"omitempty,max=64"`
}

type View struct {
	CommsXMPP      *setting.CommsXMPP `json:"comms_xmpp"`
	XMPPTargets    config.Targets     `json:"xmpp_targets"`
	SelectedTarget *config.Target     `json:"selected_target,omitempty"`
}
