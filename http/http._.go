package http

import (
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/web"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
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
) (srv *HTTP, err error) {
	srv = new(HTTP)
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

func (srv *HTTP) Startup() (err error) {
	srv.rt.Debug("status", "exec")

	if err = srv.iface.Startup(); err != nil {
		srv.rt.Error("status", "error", "error", err)
		return err
	}

	srv.rt.Info("status", "ok")

	return nil
}

func (srv *HTTP) Run() (err error) {
	srv.rt.Debug("status", "exec")

	if err = srv.iface.Run(); err != nil {
		srv.rt.Error("status", "error", "error", err)
		return err
	}

	srv.rt.Info("status", "ok")

	return nil
}

func (srv *HTTP) Shutdown() (err error) {
	srv.rt.Debug("status", "exec")

	if err = srv.iface.Shutdown(); err != nil {
		srv.rt.Error("status", "error", "error", err)
		return err
	}

	srv.rt.Info("status", "ok")

	return nil
}
