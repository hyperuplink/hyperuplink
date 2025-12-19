package topic

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Kind string

const (
	Regular  Kind = "regular"
	Question Kind = "question"
	Poll     Kind = "poll"
	RSVP     Kind = "rsvp"
)

type Topic struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`

	ForumID  uuid.UUID `json:"forum_id"`
	AuthorID uuid.UUID `json:"author_id"`

	Kind        Kind     `json:"kind"`
	Anonymous   bool     `json:"anonymous"`
	Pinned      bool     `json:"pinned"`
	Text        string   `json:"text"`
	HTML        string   `json:"html"`
	PollOptions []string `json:"poll_options"`

	Views int64 `json:"views"`

	CreatedAt   pgtype.Timestamp `json:"created_at"`
	UpdatedAt   pgtype.Timestamp `json:"updated_at"`
	EndedAt     pgtype.Timestamp `json:"ended_at"`
	ModeratedAt pgtype.Timestamp `json:"moderated_at"`
	SpammedAt   pgtype.Timestamp `json:"spammed_at"`
	LockedAt    pgtype.Timestamp `json:"locked_at"`
	DeletedAt   pgtype.Timestamp `json:"deleted_at"`
}
