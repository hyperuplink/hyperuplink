package setting

type Topics struct {
	AllowKindQuestion bool `json:"allow_kind_question"`
	AllowKindPoll     bool `json:"allow_kind_poll"`
	AllowKindRSVP     bool `json:"allow_kind_rsvp"`
}
