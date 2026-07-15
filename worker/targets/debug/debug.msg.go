package debug

import (
	"fmt"
	"strings"
)

type Msg struct {
	filename string
	body     string
}

func NewMsg() (msg *Msg) {
	msg = new(Msg)
	return msg
}

func (msg *Msg) ToString() string {
	return msg.body
}

func (msg *Msg) Filename() string {
	return msg.filename
}

type msgBuilder struct {
	sb strings.Builder
}

func (b *msgBuilder) header(key string, value string) {
	fmt.Fprintf(&b.sb, "%s: %s\n", key, value)
}

func (b *msgBuilder) section(name string, body string) {
	fmt.Fprintf(&b.sb, "\n--- %s ---\n%s", name, body)
	if !strings.HasSuffix(body, "\n") {
		b.sb.WriteString("\n")
	}
}

func (b *msgBuilder) String() string {
	return b.sb.String()
}
