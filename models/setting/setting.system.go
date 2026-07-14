package setting

const (
	DEFAULT_TOPICS_PER_PAGE          int = 10
	DEFAULT_POSTS_PER_PAGE           int = 10
	DEFAULT_ADMIN_LOG_RETENTION_DAYS int = 30
)

type System struct {
	Name                  string `json:"name"`
	BaseURL               string `json:"base_url"`
	TopicsPerPage         int    `json:"topics_per_page"`
	PostsPerPage          int    `json:"posts_per_page"`
	AdminLogRetentionDays int    `json:"admin_log_retention_days"`
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

func (sys *System) GetAdminLogRetentionDays() (days int) {
	if sys.AdminLogRetentionDays == 0 {
		sys.AdminLogRetentionDays = DEFAULT_ADMIN_LOG_RETENTION_DAYS
	}

	return sys.AdminLogRetentionDays
}
