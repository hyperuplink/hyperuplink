package vtopic

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/mrusme/hyperuplink/models/topic"
)

type VTopic struct {
	topic.Topic

	Replies     int              `json:"replies"`
	Views       int              `json:"views"`
	LastReplyAt pgtype.Timestamp `json:"last_reply_at"`
}
