package http

import (
	"github.com/mrusme/hyperuplink/errs"
	"github.com/mrusme/hyperuplink/http/web"
	"github.com/mrusme/hyperuplink/runtime"
)

type InterfaceType int

const (
	IfaceAPI InterfaceType = iota
	IfaceWeb
)

type Iface interface {
	Startup() error
	Run() error
	Shutdown() error
}

type HTTP struct {
	rt    *runtime.Runtime
	iface Iface
}

func New(
	rt *runtime.Runtime,
	ifType InterfaceType,
) (*HTTP, error) {
	var err error

	srv := new(HTTP)

	srv.rt = rt

	switch ifType {
	case IfaceAPI:
		// TODO
	case IfaceWeb:
		if srv.iface, err = web.New(srv.rt); err != nil {
			return nil, err
		}
	default:
		return nil, errs.ErrIfaceTypeUnsupported
	}

	return srv, nil
}

func (srv *HTTP) Startup() error {
	var err error

	srv.rt.Debug("status", "exec")

	if err = srv.iface.Startup(); err != nil {
		srv.rt.Error("status", "error", "error", err)
		return err
	}

	srv.rt.Info("status", "ok")

	return nil
}

func (srv *HTTP) Run() error {
	var err error

	srv.rt.Debug("status", "exec")

	if err = srv.iface.Run(); err != nil {
		srv.rt.Error("status", "error", "error", err)
		return err
	}

	srv.rt.Info("status", "ok")

	return nil
}

func (srv *HTTP) Shutdown() error {
	var err error

	srv.rt.Debug("status", "exec")

	if err = srv.iface.Shutdown(); err != nil {
		srv.rt.Error("status", "error", "error", err)
		return err
	}

	srv.rt.Info("status", "ok")

	return nil
}
