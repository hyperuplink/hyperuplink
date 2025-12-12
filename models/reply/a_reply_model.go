package reply

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Kind string

const (
	Regular Kind = "regular"
	Vote    Kind = "vote"
	RSVP    Kind = "rsvp"
)

type RSVPResponse int

const (
	Yes   RSVPResponse = 0
	No    RSVPResponse = 1
	Maybe RSVPResponse = 2
)

type Reply struct {
	ID uuid.UUID `json:"id"`

	TopicID  uuid.UUID     `json:"topic_id"`
	ReplyID  uuid.NullUUID `json:"reply_id"`
	AuthorID uuid.UUID     `json:"author_id"`

	Kind         Kind         `json:"kind"`
	Text         string       `json:"text"`
	PollVote     int          `json:"poll_vote"`
	RSVPResponse RSVPResponse `json:"rsvp_response"`

	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
	ModeratedAt pgtype.Timestamp `json:"moderated_at"`
	SpammedAt   pgtype.Timestamp `json:"spammed_at"`
	DeletedAt   pgtype.Timestamp `json:"deleted_at"`
}
