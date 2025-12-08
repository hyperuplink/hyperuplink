package signupconfirmation

import (
	"fmt"

	"github.com/mrusme/hyperuplink/models/asyncjob/common"
	"github.com/mrusme/hyperuplink/models/setting"
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
	sys *setting.System,
	forUser *user.User,
	subject string,
	token string,
	signupURL string,
) (entity SignupConfirmation, err error) {
	entity = SignupConfirmation{}
	entity.SetSystem(sys)
	entity.SetRecipient(forUser)
	entity.SetSubject(subject)
	entity.Signup.SetSignup(token,
		fmt.Sprintf("%s/%s", entity.System.BaseURL, signupURL),
	)

	return entity, nil
}

func (entity SignupConfirmation) SetSystem(sys *setting.System) {
	entity.System.SetSystem(sys)
}

func (entity SignupConfirmation) SetRecipient(rcpt *user.User) {
	entity.Recipient.SetRecipient(rcpt)
}

func (entity SignupConfirmation) SetSubject(subject string) {
	entity.Subject = subject
}

func (entity Signup) SetSignup(token string, url string) {
	entity.Token = token
	entity.URL = url
}
