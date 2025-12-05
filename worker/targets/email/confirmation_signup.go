package email

import (
	"github.com/mrusme/hyperuplink/errs"
	"github.com/mrusme/hyperuplink/models/asyncjob"
	"github.com/mrusme/hyperuplink/models/asyncjob/confirmation/signupconfirmation"
)

func (t *Email) ExecuteConfirmationSignup(
	job asyncjob.AsyncJob,
	args *Args,
	payloads []signupconfirmation.SignupConfirmation,
) (err error) {
	t.rt.Info("execute target", "email",
		"type", job.Type, "sub_type", job.SubType)

	return errs.ErrNotImplemented
}
