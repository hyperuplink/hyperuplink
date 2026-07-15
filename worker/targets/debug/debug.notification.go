package debug

import (
	"encoding/json"

	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/notification/replynotification"
)

func (t *Debug) ExecuteNotification(
	job asyncjob.AsyncJob,
	args *Args,
) (err error) {
	t.rt.Info("execute target", "debug",
		"type", job.Type, "sub_type", job.SubType)

	switch job.SubType {
	case asyncjob.Reply:
		var payload replynotification.ReplyNotification
		var payloads []replynotification.ReplyNotification
		if job.Batch == true {
			if err := json.Unmarshal(job.Payload, &payloads); err != nil {
				return err
			}
		} else {
			if err := json.Unmarshal(job.Payload, &payload); err != nil {
				return err
			}
			payloads = append(payloads, payload)
		}
		return t.ExecuteNotificationReply(job, args, payloads)
	case asyncjob.Signin, asyncjob.Mention:
		// TODO
		return errs.ErrNotImplemented
	default:
		return errs.ErrJobSubTypeInvalid
	}
}
