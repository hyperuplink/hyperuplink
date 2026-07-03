package vsearchresult

import (
	"math"

	"github.com/jackc/pgx/v5/pgtype"
)

type VSearchResult struct {
	Kind           string           `json:"kind"` // "topic" | "reply"
	CategorySlug   string           `json:"category_slug"`
	ForumSlug      string           `json:"forum_slug"`
	TopicSlug      string           `json:"topic_slug"`
	TopicName      string           `json:"topic_name"`
	AuthorUsername string           `json:"author_username"`
	CreatedAt      pgtype.Timestamp `json:"created_at"`

	// ReplyShortID is short id of matched reply (empty for topic hits).
	// Used as #post-<id> anchor on topic page.
	ReplyShortID string `json:"reply_short_id"`
	// ReplyPosition is 1-based index of matched reply within topic.
	// (Ordered by created_at, 0 for topic hits)
	ReplyPosition int64 `json:"reply_position"`
}

func (r VSearchResult) IsReply() bool {
	return r.Kind == "reply"
}

func (r VSearchResult) PageWithin(perPage int) int {
	if perPage <= 0 || !r.IsReply() {
		return 1
	}
	return int(math.Ceil(float64(r.ReplyPosition+1) / float64(perPage)))
}
