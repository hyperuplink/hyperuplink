package auth

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

type UpdateInput struct {
	AddressType int `json:"address_type" form:"address_type" validate:"oneof=0 1 2"`
}

func Update(
	rt *runtime.Runtime,
	actorID uuid.NullUUID,
	in *UpdateInput,
) (err error) {
	settingAuth, err := settingRepo.GetByID[setting.Auth](
		gh.Repositories(rt).Setting,
		"auth",
	)
	if err != nil {
		return err
	}

	before := settingAuth.JSONValue

	settingAuth.JSONValue.AddressType = setting.AddressType(in.AddressType)

	if err = settingRepo.Update[setting.Auth](
		gh.Repositories(rt).Setting,
		settingAuth,
	); err != nil {
		return err
	}

	if actorID.Valid {
		logicactivity.RecordAdminSettingsUpdate(rt, actorID.UUID,
			"auth", before, settingAuth.JSONValue)
	}

	return nil
}
