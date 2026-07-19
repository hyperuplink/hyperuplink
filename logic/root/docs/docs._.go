package docs

const (
	PageAbout   = "about"
	PageContact = "contact"
	PagePrivacy = "privacy"
	PageTerms   = "terms"
)

type Page struct {
	Page string `json:"page"`
	Text string `json:"text"`
	HTML string `json:"html"`
}
