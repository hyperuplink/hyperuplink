package activity

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Kind string

const (
	TopicView           Kind = "topic_view"
	AdminVisit          Kind = "admin_visit"
	AdminSettingsUpdate Kind = "admin_settings_update"
)

type Subject string

const (
	Topic   Subject = "topic"
	Reply   Subject = "reply"
	Setting Subject = "setting"
	Page    Subject = "page"
)

type Policy int

const (
	Buffered Policy = iota
	Immediate
)

func (k Kind) Policy() Policy {
	switch k {
	case TopicView:
		return Buffered
	default:
		return Immediate
	}
}

func (k Kind) Coalesces() bool {
	switch k {
	case TopicView:
		return true
	default:
		return false
	}
}

type Key struct {
	Kind      Kind
	ActorID   uuid.UUID
	DedupeKey string
}

type Record struct {
	Key

	Subject   Subject
	SubjectID uuid.NullUUID
	Context   []byte
	Count     int
}

type Activity struct {
	ID        uuid.UUID     `json:"id"`
	Kind      Kind          `json:"kind"`
	ActorID   uuid.UUID     `json:"actor_id"`
	Subject   Subject       `json:"subject"`
	SubjectID uuid.NullUUID `json:"subject_id"`
	DedupeKey pgtype.Text   `json:"dedupe_key"`
	Count     int           `json:"count"`
	Context   []byte        `json:"context"`

	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	DeletedAt pgtype.Timestamp `json:"deleted_at"`
}

func NewTopicView(actorID uuid.UUID, topicID uuid.UUID) Record {
	return Record{
		Key: Key{
			Kind:      TopicView,
			ActorID:   actorID,
			DedupeKey: topicID.String(),
		},
		Subject:   Topic,
		SubjectID: uuid.NullUUID{UUID: topicID, Valid: true},
		Count:     1,
	}
}
