package email

import (
	"github.com/mrusme/hyperuplink/models/asyncjob"
	"github.com/mrusme/hyperuplink/models/asyncjob/notification/replynotification"
	"github.com/wneessen/go-mail"
)

func (t *Email) ExecuteNotificationReply(
	job asyncjob.AsyncJob,
	args *Args,
	payloads []replynotification.ReplyNotification,
) (err error) {
	t.rt.Info("execute target", "email",
		"type", job.Type, "sub_type", job.SubType)

	t.rt.Debug(
		"payloads", payloads,
	)

	var messages []*mail.Msg
	for _, payload := range payloads {
		var message *mail.Msg

		if message, err = t.prepareMessage(
			job.Type,
			job.SubType,
			t.AddPlusToAddr(
				t.def.Email.From.Email,
				payload.Reply.ID.String(),
			),
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
