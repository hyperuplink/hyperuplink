package common

import (
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type Recipient struct {
	Username string
	Address  string
	Lang     string
}

func (entity *Recipient) SetRecipient(rcpt *user.User) {
	entity.Username = rcpt.Username
	entity.Address = rcpt.Email
	entity.Lang = rcpt.Language
}

type System struct {
	Name    string
	BaseURL string
}

func (entity *System) SetSystem(sys *setting.System) {
	entity.Name = sys.Name
	entity.BaseURL = sys.BaseURL
}
