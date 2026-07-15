package debug

import (
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/notification/replynotification"
)

func (t *Debug) ExecuteNotificationReply(
	job asyncjob.AsyncJob,
	args *Args,
	payloads []replynotification.ReplyNotification,
) (err error) {
	t.rt.Info("execute target", "debug",
		"type", job.Type, "sub_type", job.SubType)

	var messages []*Msg
	for _, payload := range payloads {
		var message *Msg

		if message, err = t.prepareMessage(
			job,
			payload.Recipient.Username,
			payload.Recipient.Address,
			payload.Recipient.Lang,
			payload.Subject,
			payload,
		); err != nil {
			return err
		}

		messages = append(messages, message)
	}

	return t.SendMessages(messages)
}
