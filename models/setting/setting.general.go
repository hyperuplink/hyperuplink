package setting

type General struct {
	EnablePrivacyPolicy bool   `json:"enable_privacy_policy"`
	PrivacyPolicy       string `json:"privacy_policy"`

	EnableTerms bool   `json:"enable_terms"`
	Terms       string `json:"terms"`

	EnableContact bool   `json:"enable_contact"`
	Contact       string `json:"contact"`

	EnableAbout bool   `json:"enable_about"`
	About       string `json:"about"`

	EnableQuit bool   `json:"enable_quit"`
	QuitURL    string `json:"quit_url"`
}
