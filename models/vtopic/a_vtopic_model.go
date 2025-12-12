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

	AuthorUsername       string           `json:"author_username"`
	AuthorRole           string           `json:"author_role"`
	AuthorMemberOf       []string         `json:"author_member_of"`
	AuthorEmail          string           `json:"author_email"`
	AuthorProfilePicture string           `json:"author_profile_picture"`
	AuthorSignature      string           `json:"author_signature"`
	AuthorJoinedAt       pgtype.Timestamp `json:"author_joined_at"`

	Replies     int              `json:"replies"`
	LastReplyAt pgtype.Timestamp `json:"last_reply_at"`
}
