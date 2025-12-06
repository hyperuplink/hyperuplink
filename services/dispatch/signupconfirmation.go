package dispatch

import (
	"github.com/mrusme/hyperuplink/models/asyncjob"
	"github.com/mrusme/hyperuplink/models/asyncjob/confirmation/signupconfirmation"
)

func (disp *Dispatch) SignupConfirmations(
	targetID string,
	payload []signupconfirmation.SignupConfirmation,
) (err error) {
	j := asyncjob.New(
		targetID,
		asyncjob.Confirmation,
		asyncjob.Signup,
	)
	if err = j.SetPayload(payload); err != nil {
		return err
	}

	if err = disp.Job(j); err != nil {
		return err
	}

	return nil
}

func (disp *Dispatch) SignupConfirmation(
	targetID string,
	payload signupconfirmation.SignupConfirmation,
) (err error) {
	return disp.SignupConfirmations(
		targetID,
		[]signupconfirmation.SignupConfirmation{payload},
	)
}
