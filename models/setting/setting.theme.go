package setting

type Theme struct {
	ThemeStorageProviderID string `json:"theme_storage_provider_id"`
	ThemeStoragePath       string `json:"theme_storage_path"`
	CustomBanner           string `json:"custom_banner"`
	CustomFavicon          string `json:"custom_favicon"`
}
