package site

type Breadcrumb struct {
	IsActive bool
	Label    string
	Title    string
	Href     string
}

func (s *Site) Breadcrumbs() []Breadcrumb {
	return []Breadcrumb{
		{
			Label: "Cloud",
			Title: "Cloud",
			Href:  s.HrefTo(""), // TODO: Link
		},
		{
			Label: "Hosting",
			Title: "Hosting",
			Href:  s.HrefTo(""), // TODO: Link
		},
		{
			Label: "Opinions on Clouvider?",
			Title: "Opinions on Clouvider?",
			Href:  s.HrefTo(""), // TODO: Link
			IsActive: true,
		},
	}
}
