package reply

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Reply struct {
	ID      uuid.UUID `json:"id"`
	ShortID string    `json:"short_id"`

	TopicID  uuid.UUID     `json:"topic_id"`
	ReplyID  uuid.NullUUID `json:"reply_id"`
	AuthorID uuid.UUID     `json:"author_id"`

	Text string `json:"text"`
	HTML string `json:"html"`

	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
	ModeratedAt pgtype.Timestamp `json:"moderated_at"`
	SpammedAt   pgtype.Timestamp `json:"spammed_at"`
	DeletedAt   pgtype.Timestamp `json:"deleted_at"`
}
