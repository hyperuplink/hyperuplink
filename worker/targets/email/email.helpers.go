package email

import (
	"fmt"
	htmltemplate "html/template"
	"strings"
	texttemplate "text/template"

	"github.com/wneessen/go-mail"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
)

type TmplCacheItem struct {
	TextTmpl *texttemplate.Template
	HtmlTmpl *htmltemplate.Template
}

type TmplCache map[string]TmplCacheItem

func tmplCacheKey(
	jobType asyncjob.JobType,
	jobSubType asyncjob.JobSubType,
	lang string,
) string {
	return fmt.Sprintf("%s/%s.%s",
		string(jobType), string(jobSubType), lang,
	)
}

func (t Email) TemplatesFor(
	jobType asyncjob.JobType,
	jobSubType asyncjob.JobSubType,
	lang string,
) (
	textTmpl *texttemplate.Template,
	htmlTmpl *htmltemplate.Template,
	err error,
) {
	key := tmplCacheKey(jobType, jobSubType, lang)

	if cache, ok := t.tmplCache[key]; ok {
		textTmpl = cache.TextTmpl
		htmlTmpl = cache.HtmlTmpl
	} else {
		textTmpl, htmlTmpl, err = t.LoadTemplates(
			jobType, jobSubType, lang)
		if err != nil {
			return nil, nil, err
		}

		t.tmplCache[key] = TmplCacheItem{
			TextTmpl: textTmpl,
			HtmlTmpl: htmlTmpl,
		}
	}

	return textTmpl, htmlTmpl, nil
}

func (t *Email) LoadTemplates(
	jobType asyncjob.JobType,
	jobSubType asyncjob.JobSubType,
	lang string,
) (
	textTmpl *texttemplate.Template,
	htmlTmpl *htmltemplate.Template,
	err error,
) {
	var tmpl string = fmt.Sprintf(
		"templates/email/%s/%s.%s.tmpl",
		string(jobType), string(jobSubType), lang,
	)

	if textTmpl, err = texttemplate.ParseFS(t.rt.Embeds["templates"],
		tmpl+".eml"); err != nil {
		return nil, nil, err
	}
	if htmlTmpl, err = htmltemplate.ParseFS(t.rt.Embeds["templates"],
		tmpl+".html"); err != nil {
		return nil, nil, err
	}

	return textTmpl, htmlTmpl, nil
}

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
	textTmpl, htmlTmpl, err := t.TemplatesFor(
		jobType,
		jobSubType,
		lang,
	)
	if err != nil {
		return nil, err
	}
	// Set text template
	if err = message.SetBodyTextTemplate(
		textTmpl,
		data,
	); err != nil {
		return nil, err
	}
	// Set html template
	if err = message.AddAlternativeHTMLTemplate(
		htmlTmpl,
		data,
	); err != nil {
		return nil, err
	}

	return message, nil
}

func (t *Email) SendMessages(
	messages []*mail.Msg,
) (err error) {
	if t.rt.IsDevelopmentMode() {
		t.rt.Debug(
			"pretend", "send",
			"messages", messages,
		)
		return nil
	} else {
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

		if err = client.DialAndSend(messages...); err != nil {
			return err
		}

		return nil
	}
}
