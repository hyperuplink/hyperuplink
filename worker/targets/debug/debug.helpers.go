package debug

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
)

func sanitize(s string) string {
	var sb strings.Builder
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z',
			r >= 'A' && r <= 'Z',
			r >= '0' && r <= '9',
			r == '-', r == '_', r == '.':
			sb.WriteRune(r)
		default:
			sb.WriteRune('_')
		}
	}
	return sb.String()
}

func (t *Debug) prepareMessage(
	job asyncjob.AsyncJob,
	rcptUsername string,
	rcptAddress string,
	lang string,
	subject string,
	data interface{},
) (message *Msg, err error) {
	message = NewMsg()

	item, err := t.tmplCache.TemplatesFor(
		job.Type,
		job.SubType,
		lang,
	)
	if err != nil {
		return nil, err
	}

	b := new(msgBuilder)
	b.header("Target", fmt.Sprintf("%s (emulates %s)",
		t.def.ID, t.def.Debug.Emulates))
	b.header("Job", fmt.Sprintf("%s/%s id=%s", job.Type, job.SubType, job.ID))
	b.header("To", fmt.Sprintf("%s <%s>", rcptUsername, rcptAddress))
	b.header("Lang", lang)
	b.header("Subject", subject)

	var sb strings.Builder
	if err = item.TextTmpl.Execute(&sb, data); err != nil {
		return nil, err
	}
	b.section("text", sb.String())

	if item.HtmlTmpl != nil {
		sb.Reset()
		if err = item.HtmlTmpl.Execute(&sb, data); err != nil {
			return nil, err
		}
		b.section("html", sb.String())
	}

	message.body = b.String()
	message.filename = fmt.Sprintf("%s-%s-%s.%s-%s.txt",
		time.Now().UTC().Format("20060102T150405.000"),
		job.ID,
		job.Type,
		job.SubType,
		sanitize(rcptAddress),
	)

	return message, nil
}

func (t *Debug) SendMessages(
	messages []*Msg,
) (err error) {
	for _, message := range messages {
		if t.def.Debug.Path == "" {
			t.rt.Info("message", "debug",
				"target", t.def.ID,
				"emulates", t.def.Debug.Emulates,
				"body", message.ToString())
			continue
		}

		if err = os.MkdirAll(t.def.Debug.Path, 0o700); err != nil {
			return err
		}

		path := filepath.Join(t.def.Debug.Path, message.Filename())
		if err = os.WriteFile(path, []byte(message.ToString()), 0o600); err != nil {
			return err
		}

		t.rt.Info("wrote message", "debug",
			"target", t.def.ID,
			"emulates", t.def.Debug.Emulates,
			"path", path)
	}

	return nil
}
