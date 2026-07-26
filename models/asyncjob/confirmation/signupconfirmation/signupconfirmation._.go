package signupconfirmation

import (
	"fmt"

	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/common"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type SignupConfirmation struct {
	Recipient *common.Recipient
	Subject   string
	Signup    *Signup
	System    *common.System
}

type Signup struct {
	Token string
	Path  string
	URL   string
}

func New(
	forUser *user.User,
	subject string,
	token string,
	signupPath string,
) (entity *SignupConfirmation, err error) {
	entity = new(SignupConfirmation)
	entity.Recipient = new(common.Recipient)
	entity.Signup = new(Signup)
	entity.System = new(common.System)

	entity.SetRecipient(forUser)
	entity.SetSubject(subject)
	entity.Signup.SetSignup(token, signupPath)

	return entity, nil
}

func (entity *SignupConfirmation) SetSystem(sys *setting.System) {
	entity.System.SetSystem(sys)
	entity.Signup.SetURL(entity.System.BaseURL)
}

func (entity *SignupConfirmation) SetRecipient(rcpt *user.User) {
	entity.Recipient.SetRecipient(rcpt)
}

func (entity *SignupConfirmation) SetSubject(subject string) {
	entity.Subject = subject
}

func (entity *Signup) SetSignup(token string, path string) {
	entity.Token = token
	entity.Path = path
}

func (entity *Signup) SetURL(baseURL string) {
	entity.URL = fmt.Sprintf("%s/%s", baseURL, entity.Path)
}

func (entity *SignupConfirmation) GetRecipient() *common.Recipient {
	return entity.Recipient
}

func (entity *SignupConfirmation) GetSubject() string {
	return entity.Subject
}
