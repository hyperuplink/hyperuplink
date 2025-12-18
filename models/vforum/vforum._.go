package vforum

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mrusme/hyperuplink/models/forum"
)

type VForum struct {
	forum.Forum

	CategoryName string `json:"category_name"`
	CategorySlug string `json:"category_slug"`

	Topics      int              `json:"topics"`
	Replies     int              `json:"replies"`
	LastReplyAt pgtype.Timestamp `json:"last_reply_at"`
}
