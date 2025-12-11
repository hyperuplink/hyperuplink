package vtopic

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mrusme/hyperuplink/models/topic"
)

type VTopic struct {
	topic.Topic

	CategoryName string `json:"category_name"`
	CategorySlug string `json:"category_slug"`

	ForumName string `json:"forum_name"`
	ForumSlug string `json:"forum_slug"`

	Author               string `json:"author"`
	AuthorEmail          string `json:"author_email"`
	AuthorProfilePicture string `json:"author_profile_picture"`
	AuthorSignature      string `json:"author_signature"`

	Replies int `json:"replies"`
	// Views       int              `json:"views"`
	LastReplyAt pgtype.Timestamp `json:"last_reply_at"`
}
