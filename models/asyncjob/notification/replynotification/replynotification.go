package replynotification

import (
	"github.com/google/uuid"
	"github.com/mrusme/hyperuplink/models/asyncjob/common"
	"github.com/mrusme/hyperuplink/models/setting"
	"github.com/mrusme/hyperuplink/models/user"
)

type ReplyNotification struct {
	Recipient common.Recipient
	Subject   string
	Reply     Reply
	System    common.System
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
	ID   uuid.UUID
	Name string
	URL  string
}

type Forum struct {
	Name string
	URL  string
}

type Category struct {
	Name string
	URL  string
}

func New(
	sys *setting.System,
	forUser *user.User,
	// forPost *post.Post,
	subject string,
) (entity ReplyNotification, err error) {
	entity = ReplyNotification{}
	entity.SetSystem(sys)
	entity.SetRecipient(forUser)
	entity.SetSubject(subject)
	// entity.SetReply(forPost)

	return entity, nil
}

func (entity ReplyNotification) SetRecipient(rcpt *user.User) {
	entity.Recipient.SetRecipient(rcpt)
}

func (entity ReplyNotification) SetSubject(subject string) {
	entity.Subject = subject
}

// func (entity ReplyNotification) SetReply(p *post.Post) {
// 	entity.Reply = Reply{
// 		ID:         p.TopicID,
// 		ByUsername: p.Author.Username,
// 		Text:       p.Text,
// 		HTML:       p.Text, // TODO: Convert to safeHTML
// 		URL:        p.Slug, // TODO: Build URL using System.BaseURL + / + slug
// 		Category: Category{
// 			Name: p.Category.Name,
// 			URL: p.Category.Slug, // TODO: Build URL using System.BaseURL + / + slug
// 		},
// 		Forum: Forum{
// 			Name: p.Forum.Name,
// 			URL: p.Forum.Slug, // TODO: Build URL using System.BaseURL + / + slug
// 		},
// 		Topic: Topic{
// 			Name: p.Forum.Name,
// 			URL: p.Forum.Slug, // TODO: Build URL using System.BaseURL + / + slug
// 		},
// 	}
// }

func (entity ReplyNotification) SetSystem(sys *setting.System) {
	entity.System.SetSystem(sys)
}
