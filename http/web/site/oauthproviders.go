package site

type OauthProvider struct {
	Label string
	Title string
	Href  string
	Class string
}

func (s *Site) OauthProviders() []OauthProvider {
	return []OauthProvider{
		{
			Label: "GitHub",
			Title: "GitHub",
			Href:  "github",
			Class: "github",
		},
	}
}
