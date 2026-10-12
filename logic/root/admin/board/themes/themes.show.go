package themes

import (
	"context"
	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	logicthemes "xn--gckvb8fzb.com/hyperuplink/logic/helpers/themes"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func Show(ctx context.Context, rt *runtime.Runtime) (view *View, err error) {
	settingTheme, err := settingRepo.GetByID[setting.Theme](ctx,
		gh.Repositories(rt).Setting,
		"theme",
	)
	if err != nil {
		return nil, err
	}

	view = new(View)
	view.Theme = &settingTheme.JSONValue

	if view.Themes, err = logicthemes.GetThemes(ctx, rt); err != nil {
		return nil, err
	}

	if view.Colorschemes, err = logicthemes.GetColorschemes(ctx, rt); err != nil {
		return nil, err
	}

	storages, err := rt.Config().Storages()
	if err != nil {
		return nil, err
	}

	for _, storage := range storages {
		view.StorageIDs = append(view.StorageIDs, storage.ID)
	}

	return view, nil
}
