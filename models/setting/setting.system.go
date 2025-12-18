package setting

const (
	DEFAULT_POSTS_PER_PAGE int = 10
)

type System struct {
	Name         string `json:"name"`
	BaseURL      string `json:"base_url"`
	PostsPerPage int    `json:"posts_per_page"`
}

func (sys *System) GetPostsPerPage() (ppp int) {
	if sys.PostsPerPage == 0 {
		sys.PostsPerPage = DEFAULT_POSTS_PER_PAGE
	}

	return sys.PostsPerPage
}
