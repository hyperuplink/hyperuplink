package email

import (
	texttemplate "text/template"
	htmltemplate "html/template"

	"github.com/mrusme/hyperuplink/models/asyncjob"
	"github.com/mrusme/hyperuplink/runtime"
	"github.com/mrusme/hyperuplink/services/config"
	"github.com/wneessen/go-mail"
)

type Email struct {
	rt        *runtime.Runtime
	targetCfg config.Target
}

func New(
	rt *runtime.Runtime,
	targetCfg config.Target,
) (*Email, error) {
	t := new(Email)

	t.rt = rt
	t.targetCfg = targetCfg

	return t, nil
}

func (t *Email) Load() error {
	t.rt.Info("load target", "email")
	t.rt.Debug("config", t.targetCfg)
	return nil
}

func (t *Email) Run() error {
	t.rt.Info("run target", "email")
	return nil
}

func (t *Email) Execute(
	j asyncjob.AsyncJob,
) (err error) {
	t.rt.Info("execute target", "email")

	recipients := j.TargetData["recipients"].([]string)
	subject := j.TargetData["subject"].(string)

	textTpl, err := texttemplate.ParseFS(t.rt.Embeds["templates"], "templates/email/"+string(j.Type)+".txt.tmpl")
	htmlTpl, err := htmltemplate.ParseFS(t.rt.Embeds["templates"], "templates/email/"+string(j.Type)+".html.tmpl")

	var messages []*mail.Msg
	for _, recipient := range recipients {
		message := mail.NewMsg()
		if err = message.EnvelopeFrom(t.targetCfg.Config["From"].(string)); err != nil {
			return err
		}
		if err = message.FromFormat(t.targetCfg.Config["FromName"].(string), t.targetCfg.Config["From"].(string)); err != nil {
			return err
		}
		if err = message.AddToFormat(recipient, recipient); err != nil {
			return err
		}
		message.SetMessageID()
		message.SetDate()
		message.SetBulk()
		message.Subject(subject)
		if err = message.SetBodyTextTemplate(textTpl, recipient); err != nil {
			return err
		}
		if err = message.AddAlternativeHTMLTemplate(htmlTpl, recipient); err != nil {
			return err
		}

		messages = append(messages, message)
	}

	client, err := mail.NewClient(
		t.targetCfg.Config["SMTPServer"].(string),
		mail.WithSMTPAuth(t.targetCfg.Config["SMTPAuthType"].(mail.SMTPAuthType)),
		mail.WithTLSPolicy(t.targetCfg.Config["SMTPTLSPolicy"].(mail.TLSPolicy)),
		mail.WithUsername(t.targetCfg.Config["SMTPUsername"].(string)),
		mail.WithPassword(t.targetCfg.Config["SMTPPassword"].(string)),
	)
	if err != nil {
		return err
	}

	if err = client.DialAndSend(messages...); err != nil {
		return err
	}

	return nil
}

func (t *Email) Shutdown() error {
	t.rt.Info("shutdown target", "email")
	return nil
}
