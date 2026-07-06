package xmpp

import (
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/notification/replynotification"
)

func (t *XMPP) ExecuteNotificationReply(
	job asyncjob.AsyncJob,
	args *Args,
	payloads []replynotification.ReplyNotification,
) (err error) {
	t.rt.Info("execute target", "xmpp",
		"type", job.Type, "sub_type", job.SubType)

	t.rt.Debug(
		"payloads", payloads,
	)

	var messages []*Msg
	for _, payload := range payloads {
		var message *Msg

		if message, err = t.prepareMessage(
			job.Type,
			job.SubType,
			"",
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
