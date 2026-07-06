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
	URL   string
}

func New(
	sys *setting.System,
	forUser *user.User,
	subject string,
	token string,
	signupURL string,
) (entity *SignupConfirmation, err error) {
	entity = new(SignupConfirmation)
	entity.Recipient = new(common.Recipient)
	entity.Signup = new(Signup)
	entity.System = new(common.System)

	entity.SetSystem(sys)
	entity.SetRecipient(forUser)
	entity.SetSubject(subject)
	entity.Signup.SetSignup(token,
		fmt.Sprintf("%s/%s", entity.System.BaseURL, signupURL),
	)

	return entity, nil
}

func (entity *SignupConfirmation) SetSystem(sys *setting.System) {
	entity.System.SetSystem(sys)
}

func (entity *SignupConfirmation) SetRecipient(rcpt *user.User) {
	entity.Recipient.SetRecipient(rcpt)
}

func (entity *SignupConfirmation) SetSubject(subject string) {
	entity.Subject = subject
}

func (entity *Signup) SetSignup(token string, url string) {
	entity.Token = token
	entity.URL = url
}
