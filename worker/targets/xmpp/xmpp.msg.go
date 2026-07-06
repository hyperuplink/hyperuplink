package xmpp

import (
	texttemplate "text/template"
)

type Msg struct {
	to string
}

func NewMsg() (msg *Msg) {
	msg = new(Msg)
	return msg
}

func (msg *Msg) To(to string) {
	msg.to = to
}

func (msg *Msg) SetBodyTextTemplate(tpl *texttemplate.Template, data interface{}) (err error) {
	// TODO: Implement
	return nil
}

func (msg *Msg) ToString()(string) {
	// TODO: Implement
	return ""
}
