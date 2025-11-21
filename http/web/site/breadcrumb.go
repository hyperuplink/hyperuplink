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
			Href:  "", // TODO: Link
		},
		{
			Label: "Hosting",
			Title: "Hosting",
			Href:  "", // TODO: Link
		},
		{
			Label: "Opinions on Clouvider?",
			Title: "Opinions on Clouvider?",
			Href:  "", // TODO: Link
			IsActive: true,
		},
	}
}
