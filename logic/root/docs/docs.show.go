package docs

import (
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func Show(
	rt *runtime.Runtime,
	page string,
) (view *Page, err error) {
	settingGeneral, err := settingRepo.GetByID[setting.General](
		rt.Repositories.Setting,
		"general",
	)
	if err != nil {
		return nil, err
	}
	general := settingGeneral.JSONValue

	var enabled bool
	var text string

	switch page {
	case PageAbout:
		enabled, text = general.EnableAbout, general.About
	case PageContact:
		enabled, text = general.EnableContact, general.Contact
	case PagePrivacy:
		enabled, text = general.EnablePrivacyPolicy, general.PrivacyPolicy
	case PageTerms:
		enabled, text = general.EnableTerms, general.Terms
	default:
		return nil, errs.ErrNoRows
	}

	if !enabled {
		return nil, errs.ErrNoRows
	}

	view = &Page{
		Page: page,
		Text: text,
	}
	if view.HTML, err = rt.Markdown.Convert(text); err != nil {
		return nil, err
	}

	return view, nil
}
