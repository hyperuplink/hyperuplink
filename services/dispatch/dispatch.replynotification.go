package dispatch

import (
	glidesdispatch "xn--gckvb8fzb.com/glides/services/dispatch"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/notification/replynotification"
)

func (disp *Dispatch) ReplyNotifications(
	payloads []*replynotification.ReplyNotification,
) (err error) {
	sys, err := disp.system()
	if err != nil {
		return err
	}

	for _, payload := range payloads {
		payload.SetSystem(sys)
	}

	return glidesdispatch.Batch(
		disp.Dispatch,
		asyncjob.Notification,
		asyncjob.Reply,
		payloads,
		asyncjob.IsJID[*replynotification.ReplyNotification],
	)
}

func (disp *Dispatch) ReplyNotification(
	payload *replynotification.ReplyNotification,
) (err error) {
	return disp.ReplyNotifications(
		[]*replynotification.ReplyNotification{payload},
	)
}
