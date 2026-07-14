package postevent

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostEventType string

const (
	Report       PostEventType = "report"
	PollVote     PostEventType = "pollvote"
	AnswerVote   PostEventType = "answervote"
	RSVPResponse PostEventType = "rsvpresponse"
)

type PostEventTarget string

const (
	Topic PostEventTarget = "topic"
	Reply PostEventTarget = "reply"
)

type ReportType int

const (
	Spam       ReportType = 0
	Misconduct ReportType = 1
	Illegal    ReportType = 2
)

type RSVPResponseType int

const (
	Yes   RSVPResponseType = 0
	No    RSVPResponseType = 1
	Maybe RSVPResponseType = 2
)

type PostEvent struct {
	ID        uuid.UUID       `json:"id"`
	Type      PostEventType   `json:"type"`
	AuthorID  uuid.UUID       `json:"author_id"`
	Target    PostEventTarget `json:"target"`
	TopicID   uuid.NullUUID   `json:"topic_id"`
	ReplyID   uuid.NullUUID   `json:"reply_id"`
	Selection int             `json:"selection"`

	CreatedAt pgtype.Timestamp `json:"created_at"`
	UpdatedAt pgtype.Timestamp `json:"updated_at"`
	DeletedAt pgtype.Timestamp `json:"deleted_at"`
}

func (m PostEvent) ReportTypeKey() string {
	switch ReportType(m.Selection) {
	case Spam:
		return "report_type_spam"
	case Misconduct:
		return "report_type_misconduct"
	case Illegal:
		return "report_type_illegal"
	default:
		return ""
	}
}
