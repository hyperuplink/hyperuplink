package targets

import (
	"xn--gckvb8fzb.com/glides/worker/targets/email"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/confirmation/signupconfirmation"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/notification/replynotification"
)

func registerEmail(t *email.Email) {
	t.Handle(asyncjob.Confirmation, asyncjob.Signup, email.HandlerFor(
		t, asyncjob.Message[*signupconfirmation.SignupConfirmation],
	))
	t.Handle(asyncjob.Notification, asyncjob.Reply, email.HandlerFor(
		t, asyncjob.Message[*replynotification.ReplyNotification],
	))
}
