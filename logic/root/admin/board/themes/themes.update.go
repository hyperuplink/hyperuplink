package themes

import (
	"slices"

	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
	logicthemes "xn--gckvb8fzb.com/hyperuplink/logic/helpers/themes"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func Update(
	rt *runtime.Runtime,
	actorID uuid.NullUUID,
	in *UpdateInput,
) (err error) {
	themes, err := logicthemes.GetThemes(rt)
	if err != nil {
		return err
	}
	if !slices.Contains(themes, in.Theme) {
		return errs.ErrInvalidTheme
	}

	colorschemes, err := logicthemes.GetColorschemes(rt)
	if err != nil {
		return err
	}
	if !slices.Contains(colorschemes, in.Colorscheme) {
		return errs.ErrInvalidColorscheme
	}

	storages, err := rt.Config().Storages()
	if err != nil {
		return err
	}

	valid := in.ThemeStorageProviderID == ""
	for _, storage := range storages {
		if storage.ID == in.ThemeStorageProviderID {
			valid = true
			break
		}
	}
	if !valid {
		return errs.ErrInvalidStorageProvider
	}

	settingTheme, err := settingRepo.GetByID[setting.Theme](
		gh.Repositories(rt).Setting,
		"theme",
	)
	if err != nil {
		return err
	}

	before := settingTheme.JSONValue

	settingTheme.JSONValue.Theme = in.Theme
	settingTheme.JSONValue.Colorscheme = in.Colorscheme
	settingTheme.JSONValue.ThemeStorageProviderID = in.ThemeStorageProviderID
	settingTheme.JSONValue.ThemeStoragePath = in.ThemeStoragePath

	if err = settingRepo.Update[setting.Theme](
		gh.Repositories(rt).Setting,
		settingTheme,
	); err != nil {
		return err
	}

	if actorID.Valid {
		logicactivity.RecordAdminSettingsUpdate(rt, actorID.UUID,
			"theme", before, settingTheme.JSONValue)
	}

	return nil
}
