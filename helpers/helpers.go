package helpers

import (
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/activity"
	"xn--gckvb8fzb.com/hyperuplink/services/dispatch"
	"xn--gckvb8fzb.com/hyperuplink/services/magick"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories"
)

type AuthProvider struct {
	Type   string   `koanf:"Type"`
	Key    string   `koanf:"Key"`
	Secret string   `koanf:"Secret"`
	Scopes []string `koanf:"Scopes"`
}

type AuthProviders []AuthProvider

func Activity(rt *runtime.Runtime) (service *activity.Activity) {
	if srv, _ := rt.GetService("activity"); srv != nil {
		return srv.(*activity.Activity)
	}
	return nil
}

func Dispatch(rt *runtime.Runtime) (service *dispatch.Dispatch) {
	if srv, _ := rt.GetService("dispatch"); srv != nil {
		return srv.(*dispatch.Dispatch)
	}
	return nil
}

func Magick(rt *runtime.Runtime) (service *magick.Magick) {
	if srv, _ := rt.GetService("magick"); srv != nil {
		return srv.(*magick.Magick)
	}
	return nil
}

func Repositories(rt *runtime.Runtime) (service *repositories.Repositories) {
	if srv, _ := rt.GetService("repositories"); srv != nil {
		return srv.(*repositories.Repositories)
	}
	return nil
}
