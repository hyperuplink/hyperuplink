package replynotification

import (
	"fmt"
	"html/template"

	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/models/asyncjob/common"
	"xn--gckvb8fzb.com/hyperuplink/models/reply"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type ReplyNotification struct {
	Recipient *common.Recipient
	Subject   string
	Reply     *Reply
	System    *common.System
}

type Reply struct {
	ID         uuid.UUID
	ByUsername string
	Text       string
	HTML       template.HTML
	Path       string
	URL        string
	Category   *Category
	Forum      *Forum
	Topic      *Topic
}

type Topic struct {
	ID   uuid.UUID
	Name string
	Path string
	URL  string
}

type Forum struct {
	Name string
	Path string
	URL  string
}

type Category struct {
	Name string
	Path string
	URL  string
}

func New(
	forUser *user.User,
	subject string,
) (entity *ReplyNotification, err error) {
	entity = new(ReplyNotification)
	entity.Recipient = new(common.Recipient)
	entity.Reply = new(Reply)
	entity.Reply.Category = new(Category)
	entity.Reply.Forum = new(Forum)
	entity.Reply.Topic = new(Topic)
	entity.System = new(common.System)

	entity.SetRecipient(forUser)
	entity.SetSubject(subject)

	return entity, nil
}

func (entity *ReplyNotification) SetRecipient(rcpt *user.User) {
	entity.Recipient.SetRecipient(rcpt)
}

func (entity *ReplyNotification) SetSubject(subject string) {
	entity.Subject = subject
}

func (entity *ReplyNotification) SetReply(
	rep *reply.Reply,
	byUsername string,
	path string,
) {
	entity.Reply.ID = rep.ID
	entity.Reply.ByUsername = byUsername
	entity.Reply.Text = rep.Text
	entity.Reply.HTML = template.HTML(rep.HTML)
	entity.Reply.Path = path
}

func (entity *ReplyNotification) SetCategory(name string, path string) {
	entity.Reply.Category.Name = name
	entity.Reply.Category.Path = path
}

func (entity *ReplyNotification) SetForum(name string, path string) {
	entity.Reply.Forum.Name = name
	entity.Reply.Forum.Path = path
}

func (entity *ReplyNotification) SetTopic(
	id uuid.UUID,
	name string,
	path string,
) {
	entity.Reply.Topic.ID = id
	entity.Reply.Topic.Name = name
	entity.Reply.Topic.Path = path
}

func (entity *ReplyNotification) SetSystem(sys *setting.System) {
	entity.System.SetSystem(sys)

	baseURL := entity.System.BaseURL
	entity.Reply.SetURL(baseURL)
	entity.Reply.Category.SetURL(baseURL)
	entity.Reply.Forum.SetURL(baseURL)
	entity.Reply.Topic.SetURL(baseURL)
}

func (entity *Reply) SetURL(baseURL string) {
	entity.URL = fmt.Sprintf("%s/%s", baseURL, entity.Path)
}

func (entity *Category) SetURL(baseURL string) {
	entity.URL = fmt.Sprintf("%s/%s", baseURL, entity.Path)
}

func (entity *Forum) SetURL(baseURL string) {
	entity.URL = fmt.Sprintf("%s/%s", baseURL, entity.Path)
}

func (entity *Topic) SetURL(baseURL string) {
	entity.URL = fmt.Sprintf("%s/%s", baseURL, entity.Path)
}

func (entity *ReplyNotification) GetRecipient() *common.Recipient {
	return entity.Recipient
}

func (entity *ReplyNotification) GetSubject() string {
	return entity.Subject
}
