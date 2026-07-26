package targets

import (
	glidestargets "xn--gckvb8fzb.com/glides/worker/targets"
	"xn--gckvb8fzb.com/glides/worker/targets/debug"
	"xn--gckvb8fzb.com/glides/worker/targets/email"
	"xn--gckvb8fzb.com/glides/worker/targets/xmpp"
)

func Register(id string, t glidestargets.ITarget) error {
	switch target := t.(type) {
	case *email.Email:
		registerEmail(target)
	case *xmpp.XMPP:
		registerXMPP(target)
	case *debug.Debug:
		registerDebug(target)
	}

	return nil
}
