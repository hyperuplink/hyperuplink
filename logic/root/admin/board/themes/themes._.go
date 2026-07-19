package themes

import (
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
)

type UpdateInput struct {
	Theme                  string `json:"theme" form:"theme" validate:"required"`
	Colorscheme            string `json:"colorscheme" form:"colorscheme" validate:"required"`
	ThemeStorageProviderID string `json:"theme_storage_provider_id" form:"theme_storage_provider_id" validate:"omitempty,max=64"`
	ThemeStoragePath       string `json:"theme_storage_path" form:"theme_storage_path" validate:"omitempty,max=255"`
}

type View struct {
	Theme        *setting.Theme `json:"theme"`
	Themes       []string       `json:"themes"`
	Colorschemes []string       `json:"colorschemes"`
	StorageIDs   []string       `json:"storage_ids"`
}

type ImageKind struct {
	field            func(theme *setting.Theme) *string
	errUploadFailed  error
	errAlreadySet    error
	errFormatInvalid error
}

var (
	Background = ImageKind{
		field:            func(theme *setting.Theme) *string { return &theme.CustomBackground },
		errUploadFailed:  errs.ErrBackgroundUploadFailed,
		errAlreadySet:    errs.ErrBackgroundAlreadySet,
		errFormatInvalid: errs.ErrBackgroundFormatNotAllowed,
	}
	Banner = ImageKind{
		field:            func(theme *setting.Theme) *string { return &theme.CustomBanner },
		errUploadFailed:  errs.ErrBannerUploadFailed,
		errAlreadySet:    errs.ErrBannerAlreadySet,
		errFormatInvalid: errs.ErrBannerFormatNotAllowed,
	}
	Favicon = ImageKind{
		field:            func(theme *setting.Theme) *string { return &theme.CustomFavicon },
		errUploadFailed:  errs.ErrFaviconUploadFailed,
		errAlreadySet:    errs.ErrFaviconAlreadySet,
		errFormatInvalid: errs.ErrFaviconFormatNotAllowed,
	}
)
