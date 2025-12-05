package signupconfirmation

import "github.com/mrusme/hyperuplink/models/asyncjob/common"

type SignupConfirmation struct {
	Recipient common.Recipient
	Subject   string
	Signup    Signup
}

type Signup struct {
	Token string
	URL   string
}
