package asyncjob

import (
	"github.com/google/uuid"
)

type JobType string

const (
	Confirmation JobType = "confirmation"
	Notification JobType = "notification"
)

type AsyncJob struct {
	ID         uuid.UUID              `json:"id"`
	Type       JobType                `json:"type"`
	TargetType string                 `json:"target_type"`
	TargetData map[string]interface{} `json:"target_data"`
}
