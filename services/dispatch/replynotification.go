package dispatch

import (
	"github.com/mrusme/hyperuplink/models/asyncjob"
	"github.com/mrusme/hyperuplink/models/asyncjob/notification/replynotification"
)

func (disp *Dispatch) ReplyNotifications(
	targetID string,
	payload []replynotification.ReplyNotification,
) (err error) {
	j := asyncjob.New(
		targetID,
		asyncjob.Notification,
		asyncjob.Reply,
	)
	if err = j.SetPayload(payload); err != nil {
		return err
	}

	if err = disp.Job(j); err != nil {
		return err
	}

	return nil
}

func (disp *Dispatch) ReplyNotification(
	targetID string,
	payload replynotification.ReplyNotification,
) (err error) {
	return disp.ReplyNotifications(
		targetID,
		[]replynotification.ReplyNotification{payload},
	)
}
