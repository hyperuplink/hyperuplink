package email

import (
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/config"
	"xn--gckvb8fzb.com/hyperuplink/worker/targets/tmpl"
)

type Email struct {
	rt        *runtime.Runtime
	def       config.Target
	tmplCache *tmpl.Cache
}

type Args struct{}

func New(
	rt *runtime.Runtime,
	def config.Target,
) (t *Email, err error) {
	t = new(Email)

	t.rt = rt
	t.def = def
	t.tmplCache = tmpl.NewCache(rt.Embeds["templates"], tmpl.EmailSpec)

	return t, nil
}

func (t *Email) Load() error {
	t.rt.Info("load target", "email")
	t.rt.Debug("config", t.def)
	return nil
}

func (t *Email) Run() error {
	t.rt.Info("run target", "email")
	return nil
}

func (t *Email) Shutdown() error {
	t.rt.Info("shutdown target", "email")
	return nil
}

func (t *Email) Execute(
	job asyncjob.AsyncJob,
) (err error) {
	t.rt.Info("execute target", "email")

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
