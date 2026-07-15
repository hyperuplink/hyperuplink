package dispatch

import (
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

	r, err := disp.routing()
	if err != nil {
		return err
	}

	byTarget := make(map[string][]*replynotification.ReplyNotification)
	for _, payload := range payloads {
		payload.SetSystem(sys)

		targetID := r.targetIDFor(payload.Recipient)
		byTarget[targetID] = append(byTarget[targetID], payload)
	}

	for targetID, targetPayloads := range byTarget {
		j := asyncjob.New(
			targetID,
			asyncjob.Notification,
			asyncjob.Reply,
		)
		if err = j.SetPayload(targetPayloads); err != nil {
			return err
		}

		if err = disp.Job(j); err != nil {
			return err
		}
	}

	return nil
}

func (disp *Dispatch) ReplyNotification(
	payload *replynotification.ReplyNotification,
) (err error) {
	return disp.ReplyNotifications(
		[]*replynotification.ReplyNotification{payload},
	)
}
