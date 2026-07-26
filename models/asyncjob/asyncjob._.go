package asyncjob

import (
	glidesasyncjob "xn--gckvb8fzb.com/glides/models/asyncjob"
	"xn--gckvb8fzb.com/glides/worker/targets/handler"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/common"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
)

type (
	AsyncJob   = glidesasyncjob.AsyncJob
	JobType    = glidesasyncjob.JobType
	JobSubType = glidesasyncjob.JobSubType
)

const (
	Confirmation JobType = "confirmation"
	Notification JobType = "notification"
	Cron         JobType = glidesasyncjob.Cron
)

const (
	// Confirmation
	Signup      JobSubType = "signup"
	EmailChange JobSubType = "email_change"
	// Notification
	Signin  JobSubType = "signin"
	Reply   JobSubType = "reply"
	Mention JobSubType = "mention"
	// Cron
	Run JobSubType = glidesasyncjob.Run
)

type Payload interface {
	GetRecipient() *common.Recipient
	GetSubject() string
	SetSystem(sys *setting.System)
}

func New(
	targetID string,
	jobType JobType,
	jobSubType JobSubType,
) *AsyncJob {
	return glidesasyncjob.New(targetID, jobType, jobSubType)
}

func Payloads[P any](j AsyncJob) ([]P, error) {
	return glidesasyncjob.Payloads[P](j)
}

func Message[P Payload](payload P) handler.Message {
	rcpt := payload.GetRecipient()

	return handler.Message{
		Username: rcpt.Username,
		Address:  rcpt.Address,
		Lang:     rcpt.Lang,
		Subject:  payload.GetSubject(),
	}
}

func IsJID[P Payload](payload P) bool {
	return payload.GetRecipient().IsJID
}
