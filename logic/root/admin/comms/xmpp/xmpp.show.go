package xmpp

import (
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/config"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func Show(rt *runtime.Runtime) (view *View, err error) {
	settingCommsXMPP, err := settingRepo.GetByID[setting.CommsXMPP](
		rt.Repositories.Setting,
		"comms_xmpp",
	)
	if err != nil {
		return nil, err
	}

	targets, err := rt.Config.Targets()
	if err != nil {
		return nil, err
	}

	view = new(View)
	view.CommsXMPP = &settingCommsXMPP.JSONValue
	for i := range targets {
		if !targets[i].Serves(config.TargetTypeXMPP) {
			continue
		}
		view.XMPPTargets = append(view.XMPPTargets, targets[i])
		if targets[i].ID == settingCommsXMPP.JSONValue.TargetID {
			view.SelectedTarget = &targets[i]
		}
	}

	return view, nil
}
