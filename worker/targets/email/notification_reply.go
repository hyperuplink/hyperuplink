package email

import (
	"fmt"
	"strings"

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
	messages, err = t.buildMessages(job, args, payloads)
	if err != nil {
		t.rt.Error("error", err)
		return err
	}

	if t.rt.IsDevelopmentMode() {
		t.rt.Debug(
			"pretend", "send",
			"messages", messages,
		)
		return nil
	} else {
		return t.SendMessages(messages)
	}
}

func (t *Email) buildMessages(
	job asyncjob.AsyncJob,
	args *Args,
	payloads []replynotification.ReplyNotification,
) (messages []*mail.Msg, err error) {
	for _, payload := range payloads {
		message := mail.NewMsg()

		splitFrom := strings.Split(t.def.Email.From.Email, "@")
		envFrom := fmt.Sprintf("%s+%s@%s",
			splitFrom[0],
			payload.Reply.ID.String(),
			splitFrom[1],
		)
		if err = message.EnvelopeFrom(
			envFrom,
		); err != nil {
			return messages, err
		}
		if err = message.FromFormat(
			t.def.Email.From.Name,
			t.def.Email.From.Email,
		); err != nil {
			return messages, err
		}
		if err = message.AddToFormat(
			payload.Recipient.Username,
			payload.Recipient.Address,
		); err != nil {
			return messages, err
		}
		message.SetMessageID()
		message.SetDate()
		message.SetBulk()
		message.Subject(payload.Subject)

		// Get templates
		textTmpl, htmlTmpl, err := t.TemplatesFor(
			job.Type,
			job.SubType,
			payload.Recipient.Lang,
		)
		if err != nil {
			return messages, err
		}
		// Set text template
		if err = message.SetBodyTextTemplate(
			textTmpl,
			payload,
		); err != nil {
			return messages, err
		}
		// Set html template
		if err = message.AddAlternativeHTMLTemplate(
			htmlTmpl,
			payload,
		); err != nil {
			return messages, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}
