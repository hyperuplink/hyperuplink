package signupconfirmation

import (
	"github.com/mrusme/hyperuplink/models/asyncjob/common"
	"github.com/mrusme/hyperuplink/models/user"
)

type SignupConfirmation struct {
	Recipient common.Recipient
	Subject   string
	Signup    Signup
	System    common.System
}

type Signup struct {
	Token string
	URL   string
}

func New(
	forUser *user.User,
	subject string,
	token string,
	url string,
) (entity SignupConfirmation, err error) {
	entity = SignupConfirmation{}
	entity.SetRecipient(forUser)
	entity.SetSubject(subject)
	entity.Signup.SetSignup(token, url)
	// entity.SetSystem(sys)

	return entity, nil
}

func (entity SignupConfirmation) SetRecipient(rcpt *user.User) {
	entity.Recipient.SetRecipient(rcpt)
}

func (entity SignupConfirmation) SetSubject(subject string) {
	entity.Subject = subject
}

func (entity Signup) SetSignup(token string, url string) {
	entity.Token = token
	entity.URL = url // TODO: Build URL using System.BaseURL + / + slug
}

// func (entity SignupConfirmation) SetSystem(sys *system.System) {
// 	entity.System.SetSystem(sys)
// }
