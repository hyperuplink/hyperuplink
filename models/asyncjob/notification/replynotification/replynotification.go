package replynotification

import (
	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/models/asyncjob/common"
)

type ReplyNotification struct {
	Recipient common.Recipient
	Subject   string
	Reply     Reply
}

type Reply struct {
	ID         uuid.UUID
	ByUsername string
	Text       string
	HTML       string
	URL        string
	Category   Category
	Forum      Forum
	Topic      Topic
}

type Topic struct {
	ID    uuid.UUID
	Title string
	URL   string
}

type Forum struct {
	Title string
	URL   string
}

type Category struct {
	Title string
	URL   string
}
