package xmpp

import (
	"fmt"
	htmltemplate "html/template"
	texttemplate "text/template"

	goxmpp "github.com/xmppo/go-xmpp"
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

func (t XMPP) TemplatesFor(
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

func (t *XMPP) LoadTemplates(
	jobType asyncjob.JobType,
	jobSubType asyncjob.JobSubType,
	lang string,
) (
	textTmpl *texttemplate.Template,
	htmlTmpl *htmltemplate.Template,
	err error,
) {
	var tmpl string = fmt.Sprintf(
		"templates/xmpp/%s/%s.%s.tmpl",
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

func (t *XMPP) prepareMessage(
	jobType asyncjob.JobType,
	jobSubType asyncjob.JobSubType,
	envFrom string,
	rcptUsername string,
	rcptAddress string,
	lang string,
	subject string,
	data interface{},
) (message *Msg, err error) {
	message = NewMsg()

	message.To(rcptAddress)

	// Get templates
	textTmpl, _, err := t.TemplatesFor(
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

	return message, nil
}

func (t *XMPP) SendMessages(
	messages []*Msg,
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

		for _, message := range messages {

			_, err = t.jabber.SendKeepAlive()
			if err != nil {
				t.rt.Error("failed to SendKeepAlive, attempting reconnect ...", "xmpp",
					"error", err)
				if err = t.reconnect(); err != nil {
					return err
				}
			}

			_, err = t.jabber.Send(goxmpp.Chat{
				Remote: message.to,
				Type:   "chat",
				Text:   message.ToString(),
			})
			if err != nil {
				t.rt.Error("failed to send", "xmpp",
					"error", err)
				return err
			}

			t.rt.Debug("successfully sent message", "xmpp",
				"destinationUsername", message.to)
		}
	}

	return nil
}

func (t *XMPP) reconnect() error {
	var err error

	if t.jabber != nil {
		if err = t.disconnect(); err != nil {
			return err
		}
	}

	if t.rt.IsDevelopmentMode() {
		t.rt.Debug(
			"service", "xmpp",
			"pretend", "connect",
			"host", t.jabberOpts.Host,
		)
		return nil
	} else {
		t.rt.Debug("connect to server ...", "xmpp",
			"host", t.jabberOpts.Host)
		t.jabber, err = t.jabberOpts.NewClient()
		if err != nil {
			t.rt.Error("failed to connect", "xmpp",
				"error", err)
			return err
		}
	}

	return nil
}

func (t *XMPP) disconnect() error {
	if t.rt.IsDevelopmentMode() {
		t.rt.Debug(
			"service", "xmpp",
			"pretend", "disconnect",
		)
		return nil
	} else {
		t.rt.Debug("close existing client", "xmpp")
		t.jabber.Close()
	}
	return nil
}
