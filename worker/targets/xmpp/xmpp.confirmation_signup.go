package xmpp

import (
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/confirmation/signupconfirmation"
)

func (t *XMPP) ExecuteConfirmationSignup(
	job asyncjob.AsyncJob,
	args *Args,
	payloads []signupconfirmation.SignupConfirmation,
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
