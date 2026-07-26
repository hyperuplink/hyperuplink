package email

import (
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/config"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func Show(rt *runtime.Runtime) (view *View, err error) {
	settingCommsEmail, err := settingRepo.GetByID[setting.CommsEmail](
		gh.Repositories(rt).Setting,
		"comms_email",
	)
	if err != nil {
		return nil, err
	}

	targets, err := rt.Config().Targets()
	if err != nil {
		return nil, err
	}

	view = new(View)
	view.CommsEmail = &settingCommsEmail.JSONValue
	for i := range targets {
		if !targets[i].Serves(config.TargetTypeEmail) {
			continue
		}
		view.EmailTargets = append(view.EmailTargets, targets[i])
		if targets[i].ID == settingCommsEmail.JSONValue.TargetID {
			view.SelectedTarget = &targets[i]
		}
	}

	return view, nil
}
