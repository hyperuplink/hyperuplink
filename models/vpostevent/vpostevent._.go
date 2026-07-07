package vpostevent

import (
	"fmt"
	"strings"

	"xn--gckvb8fzb.com/hyperuplink/models/postevent"
)

type VPostEvent struct {
	postevent.PostEvent

	TargetText           string `json:"target_text"`
	TargetAuthorUsername string `json:"target_author_username"`
	AuthorUsername       string `json:"author_username"`

	CategorySlug  string `json:"category_slug"`
	ForumSlug     string `json:"forum_slug"`
	TopicSlug     string `json:"topic_slug"`
	TargetShortID string `json:"target_short_id"`
}

func (m VPostEvent) TargetPath() string {
	if m.CategorySlug == "" || m.ForumSlug == "" || m.TopicSlug == "" {
		return ""
	}

	path := fmt.Sprintf("_%s/%s/%s", m.CategorySlug, m.ForumSlug, m.TopicSlug)
	if m.TargetShortID != "" {
		path = fmt.Sprintf("%s#post-%s", path, m.TargetShortID)
	}

	return path
}

func (m VPostEvent) TargetSnippet() string {
	text := strings.Join(strings.Fields(m.TargetText), " ")

	runes := []rune(text)
	if len(runes) > 32 {
		return string(runes[:32]) + "…"
	}

	return text
}
