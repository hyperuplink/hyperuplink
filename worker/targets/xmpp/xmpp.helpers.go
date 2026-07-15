package xmpp

import (
	goxmpp "github.com/xmppo/go-xmpp"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
)

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

	return message, nil
}

func (t *XMPP) SendMessages(
	messages []*Msg,
) (err error) {
	t.rt.Debug(
		"execute", "send",
		"messages", messages,
	)

	t.mu.Lock()
	defer t.mu.Unlock()

	for _, message := range messages {
		if err = t.ensureConnected(); err != nil {
			return err
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

	return nil
}

func (t *XMPP) ensureConnected() (err error) {
	if t.jabber == nil {
		t.rt.Debug("no connection, connecting ...", "xmpp",
			"host", t.jabberOpts.Host)
		return t.reconnect()
	}

	if _, err = t.jabber.SendKeepAlive(); err != nil {
		t.rt.Error("failed to SendKeepAlive, attempting reconnect ...", "xmpp",
			"error", err)
		return t.reconnect()
	}

	return nil
}

func (t *XMPP) reconnect() (err error) {
	if t.jabber != nil {
		if err = t.disconnect(); err != nil {
			return err
		}
	}

	t.rt.Debug("connect to server ...", "xmpp",
		"host", t.jabberOpts.Host)

	t.jabber, err = t.jabberOpts.NewClient()
	if err != nil {
		t.jabber = nil
		t.rt.Error("failed to connect", "xmpp",
			"error", err)
		return err
	}

	return nil
}

func (t *XMPP) disconnect() error {
	t.rt.Debug("close existing client", "xmpp")

	t.jabber.Close()
	t.jabber = nil

	return nil
}
