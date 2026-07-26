package targets

import (
	"xn--gckvb8fzb.com/glides/worker/targets/xmpp"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/confirmation/signupconfirmation"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/notification/replynotification"
)

func registerXMPP(t *xmpp.XMPP) {
	t.Handle(asyncjob.Confirmation, asyncjob.Signup, xmpp.HandlerFor(
		t, asyncjob.Message[*signupconfirmation.SignupConfirmation],
	))
	t.Handle(asyncjob.Notification, asyncjob.Reply, xmpp.HandlerFor(
		t, asyncjob.Message[*replynotification.ReplyNotification],
	))
}
