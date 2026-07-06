package email

import (
	"encoding/json"

	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/confirmation/signupconfirmation"
)

func (t *Email) ExecuteConfirmation(
	job asyncjob.AsyncJob,
	args *Args,
) (err error) {
	t.rt.Info("execute target", "email",
		"type", job.Type, "sub_type", job.SubType)
	switch job.SubType {
	case asyncjob.Signup:
		var payload signupconfirmation.SignupConfirmation
		var payloads []signupconfirmation.SignupConfirmation
		if job.Batch == true {
			if err := json.Unmarshal(job.Payload, &payloads); err != nil {
				return err
			}
		} else {
			if err := json.Unmarshal(job.Payload, &payload); err != nil {
				return err
			}
			payloads = append(payloads, payload)
		}
		return t.ExecuteConfirmationSignup(job, args, payloads)
	case asyncjob.EmailChange:
		// TODO
		return errs.ErrNotImplemented
	default:
		return errs.ErrJobSubTypeInvalid
	}
}
