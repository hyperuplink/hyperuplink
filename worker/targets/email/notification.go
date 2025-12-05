package email

import (
	"encoding/json"

	"github.com/mrusme/hyperuplink/errs"
	"github.com/mrusme/hyperuplink/models/asyncjob"
	"github.com/mrusme/hyperuplink/models/asyncjob/notification/reply"
)

func (t *Email) ExecuteNotification(
	job asyncjob.AsyncJob,
	args *Args,
) (err error) {
	t.rt.Info("execute target", "email",
		"type", job.Type)
	switch job.SubType {
	case asyncjob.Signin:
		// TODO
		return errs.ErrNotImplemented
	case asyncjob.Reply:
		var payload reply.Reply
		var payloads []reply.Reply
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
	default:
		return errs.ErrJobSubTypeInvalid
	}
}
