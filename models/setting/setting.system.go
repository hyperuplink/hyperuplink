package setting

const (
	DEFAULT_TOPICS_PER_PAGE int = 10
	DEFAULT_POSTS_PER_PAGE  int = 10

	DEFAULT_THEME       string = "macos9"
	DEFAULT_COLORSCHEME string = "hyperuplink-light"
)

type System struct {
	Name          string `json:"name"`
	BaseURL       string `json:"base_url"`
	Theme         string `json:"theme"`
	Colorscheme   string `json:"colorscheme"`
	TopicsPerPage int    `json:"topics_per_page"`
	PostsPerPage  int    `json:"posts_per_page"`
}

func (sys *System) GetTheme() (theme string) {
	if sys.Theme == "" {
		sys.Theme = DEFAULT_THEME
	}

	return sys.Theme
}

func (sys *System) GetColorscheme() (colorscheme string) {
	if sys.Colorscheme == "" {
		sys.Colorscheme = DEFAULT_COLORSCHEME
	}

	return sys.Colorscheme
}

func (sys *System) GetTopicsPerPage() (pp int) {
	if sys.TopicsPerPage == 0 {
		sys.TopicsPerPage = DEFAULT_TOPICS_PER_PAGE
	}

	return sys.TopicsPerPage
}

func (sys *System) GetPostsPerPage() (pp int) {
	if sys.PostsPerPage == 0 {
		sys.PostsPerPage = DEFAULT_POSTS_PER_PAGE
	}

	return sys.PostsPerPage
}
