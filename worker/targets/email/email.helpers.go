package email

import (
	"fmt"
	"strings"

	"github.com/wneessen/go-mail"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
)

func (t *Email) AddPlusToAddr(email string, plus string) (plusemail string) {
	splitAddr := strings.Split(email, "@")
	plusemail = fmt.Sprintf("%s+%s@%s",
		splitAddr[0],
		plus,
		splitAddr[1],
	)
	return plusemail
}

func (t *Email) prepareMessage(
	jobType asyncjob.JobType,
	jobSubType asyncjob.JobSubType,
	envFrom string,
	rcptUsername string,
	rcptAddress string,
	lang string,
	subject string,
	data interface{},
) (message *mail.Msg, err error) {
	message = mail.NewMsg()

	if err = message.EnvelopeFrom(envFrom); err != nil {
		return nil, err
	}

	if err = message.FromFormat(
		t.def.Email.From.Name,
		t.def.Email.From.Email,
	); err != nil {
		return nil, err
	}

	if err = message.AddToFormat(
		rcptUsername,
		rcptAddress,
	); err != nil {
		return nil, err
	}

	message.SetMessageID()
	message.SetDate()
	message.SetBulk()
	message.Subject(subject)

	// Get templates
	item, err := t.tmplCache.TemplatesFor(
		jobType,
		jobSubType,
		lang,
	)
	if err != nil {
		return nil, err
	}
	// Set text template
	if err = message.SetBodyTextTemplate(
		item.TextTmpl,
		data,
	); err != nil {
		return nil, err
	}
	// Set html template
	if err = message.AddAlternativeHTMLTemplate(
		item.HtmlTmpl,
		data,
	); err != nil {
		return nil, err
	}

	return message, nil
}

func (t *Email) SendMessages(
	messages []*mail.Msg,
) (err error) {
	t.rt.Debug(
		"execute", "send",
		"messages", messages,
	)

	client, err := mail.NewClient(
		t.def.Email.SMTPServer,
		mail.WithSMTPAuth(mail.SMTPAuthType(t.def.Email.SMTPAuthType)),
		mail.WithTLSPolicy(mail.TLSPolicy(t.def.Email.SMTPTLSPolicy)),
		mail.WithUsername(t.def.Email.SMTPUsername),
		mail.WithPassword(t.def.Email.SMTPPassword),
	)
	if err != nil {
		return err
	}

	return client.DialAndSend(messages...)
}
