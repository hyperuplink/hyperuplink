package xmpp

import (
	"strings"
	texttemplate "text/template"
)

type Msg struct {
	to   string
	body string
}

func NewMsg() (msg *Msg) {
	msg = new(Msg)
	return msg
}

func (msg *Msg) To(to string) {
	msg.to = to
}

func (msg *Msg) SetBodyTextTemplate(tpl *texttemplate.Template, data interface{}) (err error) {
	var sb strings.Builder
	if err = tpl.Execute(&sb, data); err != nil {
		return err
	}
	msg.body = sb.String()
	return nil
}

func (msg *Msg) ToString() string {
	return msg.body
}
