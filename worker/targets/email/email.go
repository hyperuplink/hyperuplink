package email

import (
	"github.com/mrusme/hyperuplink/models/asyncjob"
	"github.com/mrusme/hyperuplink/runtime"
	"github.com/mrusme/hyperuplink/services/config"
)

type Email struct {
	rt        *runtime.Runtime
	targetCfg config.Target
}

func New(
	rt *runtime.Runtime,
	targetCfg config.Target,
) (*Email, error) {
	t := new(Email)

	t.rt = rt
	t.targetCfg = targetCfg

	return t, nil
}

func (t *Email) Load() error {
	t.rt.Info("load target", "email")
	return nil
}

func (t *Email) Run() error {
	t.rt.Info("run target", "email")
	return nil
}

func (t *Email) Execute(
	j asyncjob.AsyncJob,
) error {
	t.rt.Info("execute target", "email")
	return nil
}

func (t *Email) Shutdown() error {
	t.rt.Info("shutdown target", "email")
	return nil
}
