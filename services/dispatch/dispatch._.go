package dispatch

import (
	glidesdispatch "xn--gckvb8fzb.com/glides/services/dispatch"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	repoSetting "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

const (
	TaskJob   string = glidesdispatch.TaskJob
	TaskCron  string = glidesdispatch.TaskCron
	QueueCron string = glidesdispatch.QueueCron
)

type Dispatch struct {
	*glidesdispatch.Dispatch

	repo *repoSetting.Repository
}

func New(
	disp *glidesdispatch.Dispatch,
	repo *repoSetting.Repository,
) (d *Dispatch, err error) {
	d = new(Dispatch)

	d.Dispatch = disp
	d.repo = repo

	d.SetResolver(d)

	return d, nil
}

func (disp *Dispatch) Routing() (r glidesdispatch.Routing, err error) {
	settingCommsEmail, err := repoSetting.GetByID[setting.CommsEmail](
		disp.repo, "comms_email")
	if err != nil {
		return r, err
	}

	settingCommsXMPP, err := repoSetting.GetByID[setting.CommsXMPP](
		disp.repo, "comms_xmpp")
	if err != nil {
		return r, err
	}

	r.EmailTargetID = settingCommsEmail.JSONValue.TargetID
	r.XMPPTargetID = settingCommsXMPP.JSONValue.TargetID

	return r, nil
}

func (disp *Dispatch) system() (sys *setting.System, err error) {
	settingSystem, err := repoSetting.GetByID[setting.System](
		disp.repo, "system")
	if err != nil {
		return nil, err
	}

	return &settingSystem.JSONValue, nil
}
