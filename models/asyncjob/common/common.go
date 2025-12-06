package common

import "github.com/mrusme/hyperuplink/models/user"

type Recipient struct {
	Username string
	Address  string
	Lang     string
}

func (entity Recipient) SetRecipient(rcpt *user.User) {
	entity.Username = rcpt.Username
	entity.Address = rcpt.Email
	entity.Lang = rcpt.Language
}

type System struct {
	Name string
}

// func (entity System) SetSystem(sys *system.System) {
// 	entity.Name = sys.Name
// }
