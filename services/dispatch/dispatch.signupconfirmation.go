package dispatch

import (
	glidesdispatch "xn--gckvb8fzb.com/glides/services/dispatch"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/confirmation/signupconfirmation"
)

func (disp *Dispatch) SignupConfirmations(
	payloads []*signupconfirmation.SignupConfirmation,
) (err error) {
	sys, err := disp.system()
	if err != nil {
		return err
	}

	for _, payload := range payloads {
		payload.SetSystem(sys)
	}

	return glidesdispatch.Batch(
		disp.Dispatch,
		asyncjob.Confirmation,
		asyncjob.Signup,
		payloads,
		asyncjob.IsJID[*signupconfirmation.SignupConfirmation],
	)
}

func (disp *Dispatch) SignupConfirmation(
	payload *signupconfirmation.SignupConfirmation,
) (err error) {
	return disp.SignupConfirmations(
		[]*signupconfirmation.SignupConfirmation{payload},
	)
}
