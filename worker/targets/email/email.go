package email

import (
	"github.com/mrusme/hyperuplink/errs"
	"github.com/mrusme/hyperuplink/models/asyncjob"
	"github.com/mrusme/hyperuplink/runtime"
	"github.com/mrusme/hyperuplink/services/config"
)

type Email struct {
	rt        *runtime.Runtime
	def       config.Target
	tmplCache TmplCache
}

type Args struct{}

func New(
	rt *runtime.Runtime,
	def config.Target,
) (t *Email, err error) {
	t = new(Email)

	t.rt = rt
	t.def = def
	t.tmplCache = make(TmplCache)

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
