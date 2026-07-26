package email

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/config"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func Update(
	rt *runtime.Runtime,
	actorID uuid.NullUUID,
	in *UpdateInput,
) (err error) {
	if in.TargetID != "" {
		targets, terr := rt.Config().Targets()
		if terr != nil {
			return terr
		}

		valid := false
		for _, target := range targets {
			if target.ID == in.TargetID &&
				target.Serves(config.TargetTypeEmail) {
				valid = true
				break
			}
		}
		if !valid {
			return errs.ErrTargetIDNotFound
		}
	}

	settingCommsEmail, err := settingRepo.GetByID[setting.CommsEmail](
		gh.Repositories(rt).Setting,
		"comms_email",
	)
	if err != nil {
		return err
	}

	before := settingCommsEmail.JSONValue

	settingCommsEmail.JSONValue.TargetID = in.TargetID

	if err = settingRepo.Update[setting.CommsEmail](
		gh.Repositories(rt).Setting,
		settingCommsEmail,
	); err != nil {
		return err
	}

	if actorID.Valid {
		logicactivity.RecordAdminSettingsUpdate(rt, actorID.UUID,
			"comms_email", before, settingCommsEmail.JSONValue)
	}

	return nil
}
