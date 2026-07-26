package targets

import (
	"xn--gckvb8fzb.com/glides/worker/targets/debug"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/confirmation/signupconfirmation"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/notification/replynotification"
)

func registerDebug(t *debug.Debug) {
	t.Handle(asyncjob.Confirmation, asyncjob.Signup, debug.HandlerFor(
		t, asyncjob.Message[*signupconfirmation.SignupConfirmation],
	))
	t.Handle(asyncjob.Notification, asyncjob.Reply, debug.HandlerFor(
		t, asyncjob.Message[*replynotification.ReplyNotification],
	))
}
