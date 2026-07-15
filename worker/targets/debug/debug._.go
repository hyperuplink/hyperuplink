package debug

import (
	"fmt"

	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/config"
	"xn--gckvb8fzb.com/hyperuplink/worker/targets/tmpl"
)

type Debug struct {
	rt        *runtime.Runtime
	def       config.Target
	tmplCache *tmpl.Cache
}

type Args struct{}

func New(
	rt *runtime.Runtime,
	def config.Target,
) (t *Debug, err error) {
	t = new(Debug)

	t.rt = rt
	t.def = def

	var spec tmpl.Spec
	switch def.Debug.Emulates {
	case config.TargetTypeEmail:
		spec = tmpl.EmailSpec
	case config.TargetTypeXMPP:
		spec = tmpl.XMPPSpec
	default:
		return nil, fmt.Errorf("%w: %s",
			errs.ErrNoSuchTargetType, def.Debug.Emulates)
	}

	t.tmplCache = tmpl.NewCache(rt.Embeds["templates"], spec)

	return t, nil
}

func (t *Debug) Load() error {
	t.rt.Info("load target", "debug",
		"emulates", t.def.Debug.Emulates)
	t.rt.Debug("config", t.def)
	return nil
}

func (t *Debug) Run() error {
	t.rt.Info("run target", "debug",
		"emulates", t.def.Debug.Emulates,
		"path", t.def.Debug.Path)
	return nil
}

func (t *Debug) Shutdown() error {
	t.rt.Info("shutdown target", "debug")
	return nil
}

func (t *Debug) Execute(
	job asyncjob.AsyncJob,
) (err error) {
	t.rt.Info("execute target", "debug")

	args := new(Args)

	switch job.Type {
	case asyncjob.Confirmation:
		return t.ExecuteConfirmation(job, args)
	case asyncjob.Notification:
		return t.ExecuteNotification(job, args)
	default:
		return errs.ErrJobTypeInvalid
	}
}
