package asyncjob

import (
	"github.com/google/uuid"
)

type JobType string

const (
	Confirmation JobType = "confirmation"
	Notification JobType = "notification"
)

type JobSubType string

const (
	// Confirmation
	Signup      JobSubType = "signup"
	EmailChange JobSubType = "email_change"
	// Notification
	Signin JobSubType = "signin"
	Reply  JobSubType = "reply"
)

type AsyncJob struct {
	ID       uuid.UUID  `json:"id"`
	Type     JobType    `json:"type"`
	SubType  JobSubType `json:"sub_type"`
	TargetID string     `json:"target_id"`
	Batch    bool       `json:"batch"`
	Payload  []byte     `json:"payload"`
}
