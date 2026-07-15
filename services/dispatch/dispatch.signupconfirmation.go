package dispatch

import (
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

	r, err := disp.routing()
	if err != nil {
		return err
	}

	byTarget := make(map[string][]*signupconfirmation.SignupConfirmation)
	for _, payload := range payloads {
		payload.SetSystem(sys)

		targetID := r.targetIDFor(payload.Recipient)
		byTarget[targetID] = append(byTarget[targetID], payload)
	}

	for targetID, targetPayloads := range byTarget {
		j := asyncjob.New(
			targetID,
			asyncjob.Confirmation,
			asyncjob.Signup,
		)
		if err = j.SetPayload(targetPayloads); err != nil {
			return err
		}

		if err = disp.Job(j); err != nil {
			return err
		}
	}

	return nil
}

func (disp *Dispatch) SignupConfirmation(
	payload *signupconfirmation.SignupConfirmation,
) (err error) {
	return disp.SignupConfirmations(
		[]*signupconfirmation.SignupConfirmation{payload},
	)
}
