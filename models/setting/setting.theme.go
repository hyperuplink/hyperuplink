package setting

const (
	DEFAULT_THEME       string = "macos9"
	DEFAULT_COLORSCHEME string = "hyperuplink-light"
)

type Theme struct {
	Theme                  string `json:"theme"`
	Colorscheme            string `json:"colorscheme"`
	ThemeStorageProviderID string `json:"theme_storage_provider_id"`
	ThemeStoragePath       string `json:"theme_storage_path"`
	CustomBanner           string `json:"custom_banner"`
	CustomFavicon          string `json:"custom_favicon"`
}

func (t *Theme) GetTheme() (theme string) {
	if t.Theme == "" {
		t.Theme = DEFAULT_THEME
	}

	return t.Theme
}

func (t *Theme) GetColorscheme() (colorscheme string) {
	if t.Colorscheme == "" {
		t.Colorscheme = DEFAULT_COLORSCHEME
	}

	return t.Colorscheme
}
